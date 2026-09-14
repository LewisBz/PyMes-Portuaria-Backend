# Quickstart: validar 001-cargo-traceability

Guía para **probar** el backend de punta a punta cuando esté implementado. Contratos: [contracts/openapi.yaml](./contracts/openapi.yaml). Modelo: [data-model.md](./data-model.md). No sustituye `tasks.md` ni el código.

## Prerrequisitos

- Docker Desktop (API + PostgreSQL 16)
- `curl` o similar
- Go 1.23+ solo si corres tests en el host: `go test ./...`

Variables típicas (Compose): `DATABASE_URL`, `JWT_SECRET` (≥ 32 bytes), `UPLOAD_DIR`, `PORT=8080`.

## Setup

Desde la raíz del repo (cuando existan `deploy/docker-compose.yml` y migraciones):

```powershell
docker compose -f deploy/docker-compose.yml up --build -d
```

Esperar healthy de Postgres y que el API aplique `migrations/*.up.sql`. Comprobar:

```powershell
curl -s http://localhost:8080/api/v1/auth/login
```

Sin body debe responder JSON de error uniforme `{"error":"...","code":400}` (no HTML, no stack).

Semilla esperada de desarrollo: organización `Demo PyME`, usuario `admin@demo.local` / password documentada en el README de implementación (no commitear secretos de producción).

## Escenarios de validación

Sustituye `$TOKEN_*` tras cada login. Base: `http://localhost:8080/api/v1`.

### 1. Auth y roles (P1)

1. Login administrador → 200 con `access_token` y `rol=administrador`.
2. Crear operador y usuario cliente (con `cliente_id` tras crear cliente) → 201.
3. Login cliente y `GET /usuarios` → **403** `{code:403}`.
4. Password incorrecto → 401, mismo formato de error.

### 2. Cliente + carga (P1)

1. Como operador: `POST /clientes` → 201.
2. `POST /cargas` con `cliente_id` y `referencia` única → 201, `estado=registrada`.
3. `GET /cargas/{id}/trazabilidad` incluye evento `carga_creada`.
4. Como cliente de **otro** cliente: `GET /cargas/{id}` → 403 o 404.

### 3. Cambio de estado y línea de tiempo (P1)

1. `POST /cargas/{id}/estados` `{"estado":"en_transito"}` → 200, `estado` actualizado.
2. Trazabilidad: evento `estado_cambiado` con `estado_anterior=registrada`, `estado_nuevo=en_transito`, `actor_id` del operador.
3. Mismo estado de nuevo → **409**.
4. Cliente: el POST de estados → **403**.

### 4. Incidencias (P2)

1. `POST /cargas/{id}/incidencias` → 201 `estado=abierta` + evento `incidencia_abierta`.
2. `POST /incidencias/{id}/cerrar` → `cerrada` + evento `incidencia_cerrada`.
3. Cliente no puede abrir incidencias → 403.

### 5. Documentos (P2)

1. Operador: multipart `POST /cargas/{id}/documentos` PDF ≤ 10 MB, `visible_cliente=true` → 201.
2. Cliente dueño: `GET .../documentos` ve el ítem; `GET /documentos/{id}/download` → 200 PDF.
3. Subir con `visible_cliente=false`; cliente download → **403**.
4. Archivo `.exe` o &gt; 10 MB → **400**.

### 6. Métricas (P3)

1. Gerente: `GET /metricas` → 200 con `cargas_activas`, `cargas_retrasadas`, `incidencias_abiertas`.
2. Cliente: `GET /metricas` → **403**.

### 7. Delimitación

No deben existir rutas ni jobs hacia DIAN, terminales o navieras (revisar OpenAPI: solo `/api/v1` listado).

## Tests automatizados (cuando existan)

```powershell
go test ./tests/contract/... ./tests/integration/... ./internal/...
```

Contrato: matriz rol × endpoint alineada a `x-required-roles` en OpenAPI. Integración: transacción de cambio de estado crea exactamente un evento.

## Resultado esperado

Todos los escenarios 1–7 en verde demuestran SC-001 a SC-005 del [spec.md](./spec.md) a nivel API, sin frontend.
