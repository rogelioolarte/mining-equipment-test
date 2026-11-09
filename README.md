# Mining Equipment Test

[English](README.md) | [Español](README.es.md)

An Angular and Go application for managing mining equipment and maintenance. Built as a Full Stack technical assessment, it tracks equipment, service intervals, maintenance history, and operating hours. The frontend and API are deployed as Vercel Services, with PostgreSQL hosted on Neon.

**Demo:** [mining-equipment-test.vercel.app](https://mining-equipment-test.vercel.app/)

## Features

- Register, view, edit, and delete equipment.
- View equipment types and their maintenance intervals.
- Record maintenance and view each equipment's service history.
- Track the hour meter, hours since the last service, and service-cycle status.

## Technology

| Area | Technologies |
| --- | --- |
| Frontend | Angular, TypeScript, and Bun |
| Backend | Go, REST API, and Docker |
| Database | PostgreSQL on Neon |
| Deployment | Vercel Services |

## Architecture

The Angular frontend consumes the Go REST API. PostgreSQL stores equipment and maintenance records, and the backend applies database migrations at startup. Vercel hosts the frontend and backend; Neon hosts the database. The repository's `vercel.json` defines both services and routes `/api/*` requests to the backend.

## Local development

- Backend setup and commands: [`backend/README.md`](backend/README.md)
- Frontend setup and commands: [`frontend/README.md`](frontend/README.md)

## Deploy to production

### Prerequisites

- A GitHub repository connected to Vercel and a Vercel account with **Services** enabled. Vercel Services are currently in beta.
- A Neon project and a database role with permission to run the application's migrations.
- Bash and the Vercel CLI installed, plus an authenticated CLI session (`vercel login`).
- Docker installed and its daemon running when using local `vercel dev` or testing the backend container. `./deploy.sh` does not invoke Docker locally; it delegates builds and deployment to Vercel.
- A local root `.env` file based on [`.env.example`](.env.example). It is ignored by Git; never commit it.

### Set up Neon

Create a Neon project and database, then copy its connection details into the root `.env`: host, port, database name, role name, and role password. Use the host without a URL scheme, keep port `5432`, and set `MINNING_DB_SSLMODE=require`. Neon must accept connections from Vercel; review the project's Neon network access settings if the backend cannot connect.

The backend applies migrations during startup. Use a database role that can create and update the application's tables, and avoid deploying multiple production instances while a migration is in progress.

### Configure Vercel

1. Import the repository into Vercel and select **Services** as the project framework.
2. Keep the repository's `vercel.json` service configuration. It defines the backend container from `backend/` and the Angular frontend from `frontend/`.
3. `vercel.json` defines an `API_URL` service binding and rewrites `/api/*` to the backend. The current Angular production bundle still uses the hard-coded URL in `frontend/src/environments/environment.ts`; the default Vercel project domain works as configured. A custom domain requires changing that source URL and rebuilding/redeploying the frontend.
4. Link this checkout to the Vercel project with `vercel link`.

The backend listens on `MINNING_HTTP_PORT` (default `8080`). The deployment script sets Vercel's container `PORT` to the same value, so the service can receive traffic on the application's configured port. Set `MINNING_HTTP_HOST=0.0.0.0` and set `MINNING_HTTP_ALLOWED_ORIGINS` to the public frontend origin in the root `.env`.

### Deploy

Review the root `.env`, then run:

```bash
./deploy.sh
```

The script validates required settings and CLIs, overwrites the supported Vercel **production** variables without displaying their values, marks only `MINNING_DB_PASSWORD` as sensitive, then runs `vercel deploy --prod`. Vercel CLI authentication and project linking must already be configured. To change a setting later, update `.env` and rerun the script; Vercel environment changes apply to new deployments. Builds run on Vercel; the script does not use local Docker.

Use plain, unquoted `KEY=VALUE` lines in `.env`; the script parses assignments as data and never sources or evaluates the file. It uploads only the documented `MINNING_*` backend settings (plus the container `PORT`) and ignores `API_URL` and unrelated keys. Blank optional values are skipped, so they do not remove an existing Vercel value; remove one separately with `vercel env rm KEY production` if desired.

## Configuration notes

| Setting | Purpose |
| --- | --- |
| `MINNING_DB_HOST`, `MINNING_DB_PORT`, `MINNING_DB_NAME`, `MINNING_DB_USER`, `MINNING_DB_PASSWORD` | Neon PostgreSQL connection. Keep the password only in the ignored `.env` and Vercel's encrypted environment settings. |
| `MINNING_DB_SSLMODE` | Use `require` for Neon TLS connections. |
| `MINNING_HTTP_HOST`, `MINNING_HTTP_PORT`, `MINNING_HTTP_ALLOWED_ORIGINS` | Backend bind address, listen port, and browser origin allowlist. |
| `MINNING_ENV`, `MINNING_LOG_FORMAT`, `MINNING_LOG_LEVEL` | Production runtime and logging configuration. |
| `MINNING_HTTP_READ_TIMEOUT`, `MINNING_HTTP_WRITE_TIMEOUT`, `MINNING_DB_MAX_CONNS`, `MINNING_DB_MIN_CONNS`, `MINNING_MIGRATION_DIR_PATH`, `MINNING_MIGRATION_PLACEHOLDERS` | Optional backend settings; blank optional values are not uploaded. |

The production Angular source currently targets `https://mining-equipment-test.vercel.app/api/v1`. That project domain works as configured. For a different custom domain, update `frontend/src/environments/environment.ts` to use the new API origin and rebuild/redeploy the frontend. The backend CORS allowlist must also include the frontend's public origin.
