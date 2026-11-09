# Frontend

Angular application for the mining equipment management UI. Vercel builds it as the `frontend` service and routes `/api/*` requests to the Go backend.

## Local development

Requirements: Bun, using the version declared in `package.json`.

```bash
bun install
bun run start
```

Open `http://localhost:4200`. The development environment calls the API at `http://localhost:8080/api/v1`; start the backend separately. To create an optimized production build locally, run `bun run build`.

## API configuration

`vercel.json` defines the `API_URL` service binding and rewrites `/api/*` to the backend, but the current Angular production bundle uses the hard-coded URL in `src/environments/environment.ts`: `https://mining-equipment-test.vercel.app/api/v1`. The default project domain works as configured. For a custom domain, change that source URL and rebuild/redeploy the frontend. The backend's allowed origins must include the browser-facing frontend origin.
