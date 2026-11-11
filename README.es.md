# Prueba de equipos mineros

[English](README.md) | [Español](README.es.md)

Aplicación Angular y Go para administrar equipos mineros, intervalos de servicio, historial de mantenimiento y horas de operación. El frontend Angular consume una API REST en Go; PostgreSQL almacena los datos. Despliega el frontend y la API en Vercel y aloja PostgreSQL en Neon.

**Demo:** [mining-equipment-test.vercel.app](https://mining-equipment-test.vercel.app/)

## Funcionalidades

- Registrar, consultar, editar y eliminar equipos.
- Consultar tipos de equipo, intervalos de mantenimiento e historial de servicio.
- Seguir el horómetro, las horas desde el último servicio y el estado del ciclo.

## Arquitectura

Vercel aloja el frontend Angular y la API en Go; Neon aloja PostgreSQL. La API ejecuta las migraciones al iniciar. `vercel.json` define ambos servicios de Vercel, proporciona el enlace de servicio `API_URL` y dirige las solicitudes `/api/*` al backend. Angular usa la ruta `/api/v1`. El archivo de entorno de producción actual de Angular apunta directamente al dominio predeterminado del proyecto.

## Desarrollo local

- Configuración y comandos del backend: [`backend/README.md`](backend/README.md)
- Configuración y comandos del frontend: [`frontend/README.md`](frontend/README.md)
- El desarrollo local del backend y frontend puede requerir [Go](https://go.dev/doc/install) y [Bun](https://bun.sh/docs/installation), respectivamente; ninguna de estas herramientas es necesaria para ejecutar `./deploy.sh`.

## Despliega tu propia instancia

### Requisitos

- Un repositorio de GitHub y una cuenta de Vercel con **Services** habilitado.
- Una [cuenta de Neon](https://neon.tech/) y un rol de base de datos con permiso para crear y actualizar las tablas de la aplicación.
- Bash, [Vercel CLI](https://vercel.com/docs/cli), un checkout local de este repositorio y un archivo `.env` en la raíz, basado en [`.env.example`](.env.example). Sigue la [guía de inicio de Vercel CLI](https://vercel.com/docs/getting-started-with-vercel) para instalarla y autenticarte con `vercel login`. Git ignora `.env`; no lo confirmes en el repositorio.

`./deploy.sh` requiere Bash y una sesión autenticada de Vercel CLI. No requiere instalar ni ejecutar Docker localmente.

### 1. Crea la base de datos en Neon

Crea un proyecto y una base de datos en Neon. Copia sus datos de conexión en el `.env` de la raíz: host (sin esquema de URL), puerto (`5432`), nombre de la base, rol y contraseña del rol. Establece `MINNING_DB_SSLMODE=require`. Consulta la [guía de conexión de Neon](https://neon.tech/docs/connect/connect-from-any-app) si necesitas encontrar esos datos. Vercel debe poder conectarse a la base de datos. El backend ejecuta las migraciones al iniciar; usa un rol con permiso para crear y actualizar tablas.

### 2. Configura el proyecto de Vercel y el backend

1. Importa tu repositorio de GitHub en Vercel, habilita **Services** y conserva la configuración de `vercel.json`. Define el backend Go desde `backend/` y el frontend Angular desde `frontend/`.
2. Desde la raíz del repositorio, autentica la CLI y vincula este checkout con tu proyecto de Vercel:

   ```bash
   vercel login
   vercel link
   ```

3. En el `.env` de la raíz, configura los datos de Neon y del backend. Establece `MINNING_HTTP_HOST=0.0.0.0` y asigna a `MINNING_HTTP_ALLOWED_ORIGINS` el **origen público del frontend**, por ejemplo `https://tu-proyecto.vercel.app` (solo el origen, sin ruta). Esta es la lista de orígenes CORS permitidos en el navegador. El backend escucha en `MINNING_HTTP_PORT` (predeterminado: `8080`); el script configura el `PORT` del contenedor de Vercel con el mismo valor.

### 3. Despliega la API y el frontend

Revisa `.env` y ejecuta desde la raíz del repositorio:

```bash
./deploy.sh
```

El script valida los ajustes y las CLI requeridos, carga los ajustes `MINNING_*` admitidos al entorno de **producción** de Vercel y luego ejecuta `vercel deploy --prod`. Marca como sensible únicamente `MINNING_DB_PASSWORD`. `vercel.json` proporciona `API_URL`; el script lo ignora, así que no tienes que configurarlo manualmente en Vercel. Vercel realiza las compilaciones y el despliegue; no se usa un daemon local de Docker. Para aplicar cambios de entorno, actualiza `.env` y vuelve a desplegar.

Usa líneas simples, sin comillas, con el formato `KEY=VALUE` en `.env`. El script analiza el archivo como datos; nunca lo carga ni lo evalúa. Sube los ajustes documentados `MINNING_*` del backend (y el `PORT` del contenedor), ignora `API_URL` y otras claves, y omite valores opcionales vacíos. Un valor vacío no elimina una variable existente en Vercel; si hace falta, elimínala por separado con `vercel env rm KEY production`.

### Dominios personalizados y rutas de la API

La URL base de la API del frontend debe usar el origen de la API seguido de `/api/v1`. `vercel.json` proporciona `API_URL` y dirige `/api/*` al backend Go. Sin embargo, `frontend/src/environments/environment.ts` actualmente tiene codificada la URL `https://mining-equipment-test.vercel.app/api/v1` para el entorno de producción Angular. El dominio predeterminado del proyecto funciona tal como está. Si usas un dominio personalizado para el frontend, actualiza ese archivo para que apunte al origen de API deseado más `/api/v1`; luego vuelve a compilar y desplegar el frontend. También configura `MINNING_HTTP_ALLOWED_ORIGINS` con el origen público del frontend en `.env` y vuelve a desplegar para permitir las solicitudes del navegador mediante CORS.

## Referencia de configuración

| Ajuste | Propósito |
| --- | --- |
| `MINNING_DB_HOST`, `MINNING_DB_PORT`, `MINNING_DB_NAME`, `MINNING_DB_USER`, `MINNING_DB_PASSWORD` | Conexión a PostgreSQL en Neon. Guarda la contraseña solo en el `.env` ignorado por Git y en las variables protegidas de Vercel. |
| `MINNING_DB_SSLMODE` | Usa `require` para las conexiones TLS con Neon. |
| `MINNING_HTTP_HOST`, `MINNING_HTTP_PORT`, `MINNING_HTTP_ALLOWED_ORIGINS` | Dirección de enlace del backend, puerto de escucha y origen permitido del frontend en el navegador. |
| `MINNING_ENV`, `MINNING_LOG_FORMAT`, `MINNING_LOG_LEVEL` | Configuración del entorno de ejecución y del registro en producción. |
| `MINNING_HTTP_READ_TIMEOUT`, `MINNING_HTTP_WRITE_TIMEOUT`, `MINNING_DB_MAX_CONNS`, `MINNING_DB_MIN_CONNS`, `MINNING_MIGRATION_DIR_PATH`, `MINNING_MIGRATION_PLACEHOLDERS` | Ajustes opcionales del backend; los valores vacíos no se cargan. |
