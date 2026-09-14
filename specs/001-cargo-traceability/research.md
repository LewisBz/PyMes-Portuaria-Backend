# Research: 001-cargo-traceability

Decisiones alineadas con spec.md, plan.md, data-model.md, contracts/openapi.yaml, tasks.md y constitution v1.0.0.

## Decision: Monolito modular (un binario), no microservicios

**Rationale**: Constitution I: un `cmd/api`, dominios en `internal/*` con repositorios por interfaz. Handlers MUST NOT importar GORM. El alcance académico/PyME no justifica varios deploys.

**Alternatives considered**: Microservicios por dominio; monolito sin paquetes (GORM en handlers).

## Decision: Go 1.23+ y Standard Go Project Layout

**Rationale**: `cmd/` entrypoints, `internal/` no exportable, `pkg/httperr` transversal. Cada dominio: `model.go`, `handler.go`, `service.go`, `repository.go`, `postgres.go`.

**Alternatives considered**: Hexagonal con `pkg/` público; DDD/event bus — overkill para v1.

## Decision: Gin como HTTP, no Fiber ni net/http crudo

**Rationale**: Middleware JWT/roles, `c.Request.Context()` hasta GORM (`WithContext`). Fiber usa su propio contexto. chi es idiomático pero más boilerplate.

**Alternatives considered**: Fiber; chi; echo.

## Decision: PostgreSQL 16 + GORM + golang-migrate (SQL, no AutoMigrate en prod)

**Rationale**: Relacional (cargas, usuarios, historial). GORM detrás de interfaces. **golang-migrate** es la fuente de verdad del esquema.

Dueño de migraciones (tasks):

- `000004_create_clientes_cargas`: `clientes` (T032) y `cargas` (T036). T036 no crea `eventos_trazabilidad`.
- `000005_create_eventos_trazabilidad`: **único** DDL de la tabla de eventos (T041).
- Incidencias/documentos en migraciones posteriores.

**Alternatives considered**: sqlc; sqlx; SQLite; AutoMigrate en producción.

## Decision: JWT + bcrypt cost ≥ 12

**Rationale**: SaaS stateless. Claims **mandatorios** (nombres de producto, no solo `sub`/`org`/`role`): `user_id`, `organizacion_id`, `rol`. Middleware valida firma/exp **y** `organizaciones.activa`; login rechaza usuario o org inactivos (401). `RequireRole` por ruta; cliente acotado a su `cliente_id`.

**ID enumeration**: `GET` de carga (o trazabilidad) fuera de tenant o de otro cliente → **404**, nunca 403.

**Alternatives considered**: sesiones Redis; Argon2; OAuth/OIDC (fuera de v1).

## Decision: Multi-tenant por columna, una org en v1

**Rationale**: Una base, `organizacion_id` en toda fila de negocio. Queries filtran el claim JWT. **v1**: una sola organización por seed/despliegue; no hay API de alta de una segunda PyME.

**Alternatives considered**: schema-per-tenant; un deploy por cliente; Stripe (billing fuera de v1).

## Decision: Documentos locales con `ObjectStore`

**Rationale**: `Put`/`Get`/`Delete` en disco (`./data/uploads`). PDF/JPEG/PNG, 10 MB. Incidencias usan campo `tipo` (no severidad).

**Alternatives considered**: S3/MinIO desde el día uno; BYTEA en Postgres.

## Decision: Trazabilidad append-only; evento `carga_creada` desde US3

**Rationale**: Tabla `eventos_trazabilidad` insert-only. Tipos: `carga_creada`, `estado_cambiado`, `incidencia_abierta`, `incidencia_cerrada`, `documento_adjunto`. El alta de carga (US2) **no** escribe eventos hasta existir `000005`. Snapshot de estado en `cargas`; `fecha_entrega` se setea al pasar a `entregada`.

**Alternatives considered**: Triggers SQL; event sourcing/CQRS.

## Decision: Enum SQL inmutable de estados; `retrasada` derivado

**Rationale**: v1 no tiene catálogo editable. Enum: `registrada`, `en_transito`, `en_puerto`, `entregada`, `cancelada`. `retrasada` **no** es estado.

Flag (no persistido): `fecha_comprometida IS NOT NULL AND fecha_comprometida < NOW() AND estado != 'entregada'`.

Transiciones: no-terminal → otro distinto; mismo estado → 409; terminales no se reabren en v1.

**Alternatives considered**: FSM aduanera; strings libres; `retrasada` como valor del enum.

## Decision: Métricas SQL; tiempo promedio total de entrega

**Rationale**: `GET /api/v1/metricas` (administrador, gerente; cliente 403). Agregados: `cargas_activas`, `cargas_retrasadas` (flag derivado), `incidencias_abiertas`, `tiempo_promedio_total_entrega_horas` = `AVG(fecha_entrega - created_at)` **solo** `estado = 'entregada'`. v1 **no** calcula intervalos intermedios entre estados.

Cargas en API: create/list/get; sin PUT/DELETE destructivos (SC-004).

**Alternatives considered**: Grafana; promedio por tramo de estado; warehouse.

## Decision: Errores `{"error","code"}` en `pkg/httperr`

**Rationale**: Constitution V. 404 para recurso fuera de alcance del cliente; 403 para rol que no puede la operación (p. ej. cliente creando cargas o listando usuarios).

**Alternatives considered**: RFC 7807.

## Decision: Docker multi-stage + Compose

**Rationale**: Binario estático, Compose api + postgres + volume uploads. `context.Context` desde Gin hasta GORM.

**Alternatives considered**: Solo Go en el host.

## Decision: Contract tests primero (constitution III)

**Rationale**: `tests/contract/` en rojo antes de handlers; `tests/integration/` para schema, JWT y AuthZ (incl. gerente 200 en trazabilidad).

**Alternatives considered**: Solo tests manuales; TDD con aprobación humana de cada archivo de test.

## Resolved clarifications

| Tema | Resolución |
|------|------------|
| Frontend | Fuera de este repo |
| Tenant v1 | Columna `organizacion_id`; **una** org por seed |
| JWT | Claims `user_id`, `organizacion_id`, `rol`; org inactiva → 401 |
| HTTP | Gin + `/api/v1`; health `/healthz` operador |
| Archivos | ObjectStore + disco |
| Estados | Enum SQL inmutable; `retrasada` derivado |
| Cargas API | create/list/get |
| Trazabilidad DDL | Solo `000005_create_eventos_trazabilidad` |
| ID leak | GET carga ajena → 404 |
| Métricas | Tiempo promedio total de entrega, no tramos |
| Incidencias | Solo `tipo` |
| Billing / DIAN | Fuera de alcance |
| Governance | Constitution v1.0.0 |
