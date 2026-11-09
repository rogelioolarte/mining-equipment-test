# Prueba de equipos mineros

[English](README.md) | [Español](README.es.md)

Aplicación Angular y Go para administrar equipos mineros y su mantenimiento. Desarrollada como prueba técnica Full Stack, permite llevar el control de los equipos, los intervalos de servicio, el historial de mantenimiento y las horas de operación. El frontend y la API se despliegan como Vercel Services, y PostgreSQL está alojado en Neon.

**Demo:** [mining-equipment-test.vercel.app](https://mining-equipment-test.vercel.app/)

## Funcionalidades

- Registrar, consultar, editar y eliminar equipos.
- Consultar tipos de equipo y sus intervalos de mantenimiento.
- Registrar mantenimientos y consultar el historial de servicio de cada equipo.
- Seguir el horómetro, las horas desde el último servicio y el estado del ciclo de servicio.

## Tecnologías

| Área | Tecnologías |
| --- | --- |
| Frontend | Angular, TypeScript y Bun |
| Backend | Go, API REST y Docker |
| Base de datos | PostgreSQL en Neon |
| Despliegue | Vercel Services |

## Arquitectura

El frontend Angular consume la API REST en Go. PostgreSQL almacena los equipos y sus mantenimientos, y el backend aplica las migraciones de base de datos al iniciar. Vercel aloja el frontend y el backend; Neon aloja la base de datos. El archivo `vercel.json` del repositorio define ambos servicios y dirige las solicitudes `/api/*` al backend.

## Desarrollo local

- Configuración y comandos del backend: [`backend/README.md`](backend/README.md)
- Configuración y comandos del frontend: [`frontend/README.md`](frontend/README.md)

## Despliegue en producción

### Requisitos previos

- Un repositorio de GitHub conectado a Vercel y una cuenta de Vercel con **Services** habilitado. Vercel Services se encuentra actualmente en beta.
- Un proyecto de Neon y un rol de base de datos con permiso para ejecutar las migraciones de la aplicación.
- Bash y Vercel CLI instalados, además de una sesión de CLI autenticada (`vercel login`).
- Docker instalado y su daemon en ejecución para usar `vercel dev` localmente o probar el contenedor del backend. `./deploy.sh` no ejecuta Docker localmente: delega las compilaciones y el despliegue en Vercel.
- Un archivo `.env` en la raíz local, basado en [`.env.example`](.env.example). Git lo ignora; nunca lo confirmes en el repositorio.

### Configurar Neon

Crea un proyecto y una base de datos en Neon; luego copia los datos de conexión al archivo `.env` de la raíz: host, puerto, nombre de la base de datos, nombre del rol y contraseña del rol. Usa el host sin esquema de URL, mantén el puerto `5432` y establece `MINNING_DB_SSLMODE=require`. Neon debe aceptar conexiones desde Vercel; si el backend no puede conectarse, revisa la configuración de acceso de red del proyecto de Neon.

El backend aplica las migraciones durante el inicio. Usa un rol de base de datos que pueda crear y actualizar las tablas de la aplicación y evita desplegar varias instancias de producción mientras haya una migración en curso.

### Configurar Vercel

1. Importa el repositorio en Vercel y selecciona **Services** como framework del proyecto.
2. Conserva la configuración de servicios del `vercel.json` del repositorio. Define el contenedor del backend a partir de `backend/` y el frontend Angular a partir de `frontend/`.
3. `vercel.json` define un enlace de servicio `API_URL` y reescribe `/api/*` para dirigirlo al backend. El paquete de producción actual de Angular todavía usa la URL codificada en `frontend/src/environments/environment.ts`; el dominio predeterminado del proyecto de Vercel funciona tal como está configurado. Para usar un dominio personalizado, hay que cambiar esa URL de origen y volver a compilar y desplegar el frontend.
4. Vincula este checkout con el proyecto de Vercel mediante `vercel link`.

El backend escucha en `MINNING_HTTP_PORT` (valor predeterminado: `8080`). El script de despliegue establece el `PORT` del contenedor de Vercel con el mismo valor, para que el servicio reciba tráfico en el puerto configurado por la aplicación. En el archivo `.env` de la raíz, establece `MINNING_HTTP_HOST=0.0.0.0` y asigna a `MINNING_HTTP_ALLOWED_ORIGINS` el origen público del frontend.

### Desplegar

Revisa el archivo `.env` de la raíz y ejecuta:

```bash
./deploy.sh
```

El script valida los ajustes y las CLI requeridos, sobrescribe las variables de **producción** admitidas en Vercel sin mostrar sus valores, marca como sensible únicamente `MINNING_DB_PASSWORD` y luego ejecuta `vercel deploy --prod`. La autenticación de Vercel CLI y la vinculación del proyecto deben estar configuradas previamente. Para cambiar un ajuste, actualiza `.env` y vuelve a ejecutar el script; los cambios del entorno de Vercel se aplican a los nuevos despliegues. Las compilaciones se ejecutan en Vercel; el script no usa Docker local.

Usa líneas simples, sin comillas, con el formato `KEY=VALUE` en `.env`; el script analiza las asignaciones como datos y nunca carga ni evalúa el archivo. Solo carga los ajustes documentados del backend `MINNING_*` (más el `PORT` del contenedor) e ignora `API_URL` y las claves no relacionadas. Se omiten los valores opcionales vacíos, por lo que no eliminan un valor existente en Vercel; si quieres eliminar uno, hazlo por separado con `vercel env rm KEY production`.

## Notas de configuración

| Ajuste | Propósito |
| --- | --- |
| `MINNING_DB_HOST`, `MINNING_DB_PORT`, `MINNING_DB_NAME`, `MINNING_DB_USER`, `MINNING_DB_PASSWORD` | Conexión a PostgreSQL en Neon. Guarda la contraseña solo en el `.env` ignorado por Git y en los ajustes de entorno cifrados de Vercel. |
| `MINNING_DB_SSLMODE` | Usa `require` para las conexiones TLS con Neon. |
| `MINNING_HTTP_HOST`, `MINNING_HTTP_PORT`, `MINNING_HTTP_ALLOWED_ORIGINS` | Dirección de enlace del backend, puerto de escucha y lista de orígenes permitidos por el navegador. |
| `MINNING_ENV`, `MINNING_LOG_FORMAT`, `MINNING_LOG_LEVEL` | Configuración del entorno de ejecución y del registro en producción. |
| `MINNING_HTTP_READ_TIMEOUT`, `MINNING_HTTP_WRITE_TIMEOUT`, `MINNING_DB_MAX_CONNS`, `MINNING_DB_MIN_CONNS`, `MINNING_MIGRATION_DIR_PATH`, `MINNING_MIGRATION_PLACEHOLDERS` | Ajustes opcionales del backend; los valores opcionales vacíos no se cargan. |

El código fuente de producción de Angular actualmente apunta a `https://mining-equipment-test.vercel.app/api/v1`. Ese dominio del proyecto funciona tal como está configurado. Para usar otro dominio personalizado, actualiza `frontend/src/environments/environment.ts` con el nuevo origen de la API y vuelve a compilar y desplegar el frontend. La lista CORS permitida por el backend también debe incluir el origen público del frontend.
