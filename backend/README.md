# Backend

Go HTTP API for equipment, equipment types, and maintenance records. It connects to PostgreSQL and applies the SQL migrations in `migrations/` during startup. In production, Vercel runs the `backend` container service and sends `/api/*` requests to it.

## Local development

Requirements: Go version declared in `go.mod` and a PostgreSQL database. From this directory, create a local `.env` from the repository template and set the database values:

```bash
cp ../.env.example .env
```

Then start the API:

```bash
go run ./cmd/api
```

The default listen port is `8080`. The process loads environment variables from `.env`, validates configuration, connects to the database, applies migrations, and starts the HTTP server. Keep real credentials in local ignored files or a secret manager; do not commit `.env`.

## Configuration

Runtime settings use the `MINNING_` prefix. The repository-level [`.env.example`](../.env.example) lists available variables. A production deployment needs database host, port, name, user, and password; `MINNING_DB_SSLMODE=require` for Neon; a production environment name; and a non-empty `MINNING_HTTP_ALLOWED_ORIGINS` allowlist. Bind to `0.0.0.0` in a container.

Vercel's container `PORT` must match `MINNING_HTTP_PORT` because the Go server reads the latter (default `8080`). The deployment script sets `PORT` to that configured value. The database must be reachable from Vercel, and the configured database role must be able to apply migrations.
