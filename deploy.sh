#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$ROOT_DIR/.env"

fail() {
  printf 'Deployment setup error: %s\n' "$1" >&2
  exit 1
}

[[ -f "$ENV_FILE" ]] || fail "root .env not found; copy .env.example and fill in production settings"
command -v vercel >/dev/null 2>&1 || fail "Vercel CLI is required"
command -v bash >/dev/null 2>&1 || fail "Bash is required"
[[ -f "$ROOT_DIR/.vercel/project.json" ]] || fail "Vercel project is not linked; run vercel link first"

declare -A values=()
while IFS= read -r line || [[ -n "$line" ]]; do
  line="${line%$'\r'}"
  # Accept simple KEY=VALUE assignments only. Never evaluate shell syntax.
  [[ "$line" =~ ^[[:space:]]*$ ]] && continue
  [[ "$line" =~ ^[[:space:]]*# ]] && continue
  [[ "$line" == *=* ]] || fail "invalid .env line; expected KEY=VALUE"

  key="${line%%=*}"
  value="${line#*=}"
  case "$key" in
    MINNING_ENV|MINNING_HTTP_HOST|MINNING_HTTP_PORT|MINNING_HTTP_ALLOWED_ORIGINS|MINNING_HTTP_READ_TIMEOUT|MINNING_HTTP_WRITE_TIMEOUT|MINNING_DB_HOST|MINNING_DB_PORT|MINNING_DB_NAME|MINNING_DB_USER|MINNING_DB_PASSWORD|MINNING_DB_SSLMODE|MINNING_DB_MAX_CONNS|MINNING_DB_MIN_CONNS|MINNING_LOG_FORMAT|MINNING_LOG_LEVEL|MINNING_MIGRATION_DIR_PATH|MINNING_MIGRATION_PLACEHOLDERS)
      [[ ! -v "values[$key]" ]] || fail "duplicate .env key: $key"
      values["$key"]="$value"
      ;;
    *)
      # API_URL is supplied by the Vercel service binding; unrelated keys are ignored.
      ;;
  esac
done < "$ENV_FILE"

required=(MINNING_DB_HOST MINNING_DB_PORT MINNING_DB_NAME MINNING_DB_USER MINNING_DB_PASSWORD MINNING_HTTP_ALLOWED_ORIGINS MINNING_HTTP_PORT)
for key in "${required[@]}"; do
  [[ -n "${values[$key]:-}" ]] || fail "required .env value is missing: $key"
done

[[ "${values[MINNING_DB_SSLMODE]:-require}" == require ]] || fail "MINNING_DB_SSLMODE must be require for Neon"
[[ "${values[MINNING_HTTP_HOST]:-0.0.0.0}" == 0.0.0.0 ]] || fail "MINNING_HTTP_HOST must be 0.0.0.0 for the Vercel container"
[[ "${values[MINNING_HTTP_PORT]}" =~ ^[0-9]{1,5}$ ]] || fail "MINNING_HTTP_PORT must be a numeric port"
port_number=$((10#${values[MINNING_HTTP_PORT]}))
(( port_number >= 1 && port_number <= 65535 )) || fail "MINNING_HTTP_PORT must be between 1 and 65535"

# Production-only values are set here so a development .env cannot deploy a dev API.
values[MINNING_ENV]=production
values[MINNING_HTTP_HOST]=0.0.0.0
values[MINNING_DB_SSLMODE]=require
# Vercel container routing uses PORT; the Go application listens on MINNING_HTTP_PORT.
values[PORT]="${values[MINNING_HTTP_PORT]}"

upload_keys=(
  MINNING_ENV MINNING_HTTP_HOST MINNING_HTTP_PORT PORT MINNING_HTTP_ALLOWED_ORIGINS
  MINNING_HTTP_READ_TIMEOUT MINNING_HTTP_WRITE_TIMEOUT
  MINNING_DB_HOST MINNING_DB_PORT MINNING_DB_NAME MINNING_DB_USER MINNING_DB_PASSWORD
  MINNING_DB_SSLMODE MINNING_DB_MAX_CONNS MINNING_DB_MIN_CONNS
  MINNING_LOG_FORMAT MINNING_LOG_LEVEL MINNING_MIGRATION_DIR_PATH MINNING_MIGRATION_PLACEHOLDERS
)

for key in "${upload_keys[@]}"; do
  value="${values[$key]:-}"
  [[ -n "$value" ]] || continue
  env_args=(env add "$key" production --force --yes)
  if [[ "$key" == MINNING_DB_PASSWORD ]]; then
    env_args+=(--sensitive)
  fi
  if ! printf '%s' "$value" | vercel "${env_args[@]}" >/dev/null 2>&1; then
    fail "could not upload production variable: $key"
  fi
done

printf 'Production environment variables uploaded. Deploying Vercel project...\n'
vercel deploy --prod --yes
