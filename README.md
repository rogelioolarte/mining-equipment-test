# Mining Equipment Test

[English](README.md) | [Español](README.es.md)

An Angular and Go application for managing mining equipment, service intervals, maintenance history, and operating hours. The Angular frontend calls a Go REST API; PostgreSQL stores application data. Deploy the frontend and API to Vercel and host PostgreSQL on Neon.

**Demo:** [mining-equipment-test.vercel.app](https://mining-equipment-test.vercel.app/)

## Features

- Register, view, edit, and delete equipment.
- View equipment types, maintenance intervals, and service history.
- Track hour meters, hours since service, and service-cycle status.

## Architecture

Vercel hosts the Angular frontend and Go API; Neon hosts PostgreSQL. The API applies database migrations at startup. `vercel.json` defines both Vercel services, provides an `API_URL` service binding, and routes `/api/*` requests to the backend. Angular requests use the `/api/v1` path. The current production environment file uses the default project domain directly.

## Local development

- Backend setup and commands: [`backend/README.md`](backend/README.md)
- Frontend setup and commands: [`frontend/README.md`](frontend/README.md)
- Local backend and frontend development may require [Go](https://go.dev/doc/install) and [Bun](https://bun.sh/docs/installation), respectively; neither is required to run `./deploy.sh`.

## Deploy your own instance

### Requirements

- A GitHub repository and a Vercel account with **Services** enabled.
- A [Neon account](https://neon.tech/) and a database role allowed to create and update application tables.
- Bash, [Vercel CLI](https://vercel.com/docs/cli), a local checkout of this repository, and a root `.env` file based on [`.env.example`](.env.example). Follow the [Vercel CLI getting-started guide](https://vercel.com/docs/getting-started-with-vercel) to install the CLI and authenticate with `vercel login`. The `.env` file is ignored by Git; do not commit it.

`./deploy.sh` requires Bash and an authenticated Vercel CLI session. It does not require Docker installed or running locally.

### 1. Create the Neon database

Create a Neon project and database. Copy its connection settings into the root `.env`: host (without a URL scheme), port (`5432`), database name, role, and role password. Set `MINNING_DB_SSLMODE=require`. Check [Neon's connection guide](https://neon.tech/docs/connect/connect-from-any-app) if you need help finding the connection details. Vercel must be able to connect to the database. The backend runs migrations at startup, so use a role with permission to create and update tables.

### 2. Configure the Vercel project and backend

1. Import your GitHub repository into Vercel, enable **Services**, and keep the repository's `vercel.json` configuration. It defines the Go backend from `backend/` and the Angular frontend from `frontend/`.
2. From the repository root, authenticate and link this checkout to your Vercel project:

   ```bash
   vercel login
   vercel link
   ```

3. In the root `.env`, set the Neon values and backend settings. Set `MINNING_HTTP_HOST=0.0.0.0`, and set `MINNING_HTTP_ALLOWED_ORIGINS` to the **public frontend origin**, such as `https://your-project.vercel.app` (origin only, without a path). This is the browser CORS allowlist. The backend listens on `MINNING_HTTP_PORT` (default `8080`); the deploy script sets the Vercel container `PORT` to match.

### 3. Deploy the API and frontend

Review `.env`, then run from the repository root:

```bash
./deploy.sh
```

The script validates its required settings and CLIs, loads supported `MINNING_*` settings into Vercel's **production** environment, then runs `vercel deploy --prod`. It marks only `MINNING_DB_PASSWORD` as sensitive. `API_URL` is supplied by `vercel.json`; the script ignores it, so do not configure it manually in Vercel. Vercel performs the builds and deployment; no local Docker daemon is used. To apply environment changes, update `.env` and deploy again.

Use plain, unquoted `KEY=VALUE` lines in `.env`. The script parses the file as data and never sources or evaluates it. It uploads documented backend `MINNING_*` settings (and the container `PORT`), ignores `API_URL` and unrelated keys, and skips blank optional values. A blank value does not remove an existing Vercel variable; remove one separately with `vercel env rm KEY production` if needed.

### Custom domains and API routing

The frontend's API base URL must use the API origin followed by `/api/v1`. `vercel.json` provides `API_URL` and routes `/api/*` to the Go backend. However, `frontend/src/environments/environment.ts` currently hardcodes `https://mining-equipment-test.vercel.app/api/v1` in the Angular production environment. The default project domain works as configured. If you use a custom frontend domain, update that file to point to the intended API origin plus `/api/v1`, then rebuild and redeploy the frontend. Also set `MINNING_HTTP_ALLOWED_ORIGINS` to the frontend's public origin in `.env` and redeploy so browser requests pass CORS.

## Configuration reference

| Setting | Purpose |
| --- | --- |
| `MINNING_DB_HOST`, `MINNING_DB_PORT`, `MINNING_DB_NAME`, `MINNING_DB_USER`, `MINNING_DB_PASSWORD` | Neon PostgreSQL connection. Keep the password only in the ignored `.env` and Vercel's protected environment settings. |
| `MINNING_DB_SSLMODE` | Use `require` for Neon TLS connections. |
| `MINNING_HTTP_HOST`, `MINNING_HTTP_PORT`, `MINNING_HTTP_ALLOWED_ORIGINS` | Backend bind address, listen port, and allowed frontend browser origin. |
| `MINNING_ENV`, `MINNING_LOG_FORMAT`, `MINNING_LOG_LEVEL` | Production runtime and logging configuration. |
| `MINNING_HTTP_READ_TIMEOUT`, `MINNING_HTTP_WRITE_TIMEOUT`, `MINNING_DB_MAX_CONNS`, `MINNING_DB_MIN_CONNS`, `MINNING_MIGRATION_DIR_PATH`, `MINNING_MIGRATION_PLACEHOLDERS` | Optional backend settings; blank optional values are not uploaded. |
