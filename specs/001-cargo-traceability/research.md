# Research: 001-cargo-traceability

## Decision: Monolito modular (un binario), no microservicios

**Rationale**: El equipo y el despliegue académico/PyME no justifican red, tracing distribuido ni varios deploys. Los dominios viven en paquetes `internal/*` con límites claros (interfaces), mismo proceso.

**Alternatives considered**: Microservicios por dominio (infra excesiva); monolito sin paquetes (acoplamiento GORM en handlers).

## Decision: Go 1.23+ y Standard Go Project Layout

**Rationale**: `cmd/` para entrypoints, `internal/` no exportable, `pkg/` solo utilidades transversales (`httperr`). Coincide con la delimitación del trabajo de grado.

**Alternatives considered**: Hexagonal con `pkg/` público; DDD pesado (agregados/event bus) — overkill para v1.

## Decision: Gin como HTTP, no Fiber ni net/http crudo

**Rationale**: Middleware JWT/roles maduro, `c.Request.Context()` llega a GORM, ecosistema de binding JSON. Fiber usa su propio contexto (peor encaje con `context.Context` pedido). `net/http` + chi es idiomático pero más boilerplate para un equipo mixto.

**Alternatives considered**: Fiber (rendimiento; peor stdlib); chi (más idiomático; más código a mano); echo.

## Decision: PostgreSQL 16 + GORM + golang-migrate (SQL, no AutoMigrate en prod)

**Rationale**: Datos relacionales (cargas, usuarios, historial). GORM mapea structs y relaciones. **golang-migrate** versiona esquema con `.up.sql`/`.down.sql` para control profesional. AutoMigrate de GORM queda solo como ayuda opcional en desarrollo local, nunca como fuente de verdad en despliegue.

**Alternatives considered**: sqlc (más control, más curva); sqlx; SQLite (no escala multi-tenant ni concurrente como Postgres).

## Decision: JWT (golang-jwt/jwt/v5) + bcrypt cost 12

**Rationale**: SaaS stateless; claims `sub` (user_id), `org` (organizacion_id), `role`. Middleware global valida firma/exp; `RequireRole` y filtros por recurso (cliente → `cliente_id` de su usuario). bcrypt cost 12 equilibra seguridad y latencia de login en CPU de contenedor pequeño.

**Alternatives considered**: sesiones en Redis (estado extra); Argon2 (mejor, menos ejemplos en el curso); OAuth/OIDC (fuera de alcance v1).

## Decision: Multi-tenant por columna `organizacion_id`

**Rationale**: Una base, aislamiento lógico, encaja con “suscripción SaaS” sin schema-per-tenant. Toda query de dominio filtra tenant. Seed de una organización en desarrollo.

**Alternatives considered**: Un schema por PyME (ops complejas); un deploy por cliente (contradice SaaS único). Facturación Stripe: **fuera de v1**.

## Decision: Almacenamiento de documentos local con interfaz

**Rationale**: `ObjectStore` (`Put`, `Get`, `Delete`) implementado con disco (`./data/uploads`) en Docker volume. No acoplar el módulo `documentos` a S3. Tipos: PDF/JPEG/PNG, 10 MB.

**Alternatives considered**: S3/MinIO desde el día uno (cuenta cloud y secretos extra para la tesis); guardar BYTEA en Postgres (infla backups).

## Decision: Trazabilidad como log de eventos, no update in-place

**Rationale**: Cumple “historial de cambios, responsables e incidencias”. Tabla `eventos_trazabilidad` insert-only. Tipos: `estado_cambiado`, `incidencia_abierta`, `incidencia_cerrada`, `documento_adjunto`. El estado actual de `cargas` es un snapshot; la verdad histórica es el log.

**Alternatives considered**: Triggers SQL (menos portable a tests); Event Sourcing completo (replay, CQRS) — complejidad injustificada.

## Decision: Catálogo de estados de carga + fecha comprometida para retrasos

**Rationale**: El spec no fija la máquina de estados. v1 usa catálogo: `registrada`, `en_transito`, `en_puerto`, `entregada`, `cancelada`, más flag derivado **retrasada** si `fecha_comprometida < now` y estado no terminal. Transiciones: cualquier no-terminal → otro (salvo mismo estado, rechazado). No se modela BPMN.

**Alternatives considered**: FSM estricta puerto-real (no hay integración ni normativa DIAN); estados libres string (inconsistente para métricas).

## Decision: Métricas en SQL agregado, no warehouse

**Rationale**: Escala PyME. Endpoint `GET /api/v1/metricas` con consultas COUNT/AVG. Roles: administrador y gerente.

**Alternatives considered**: Grafana/Prometheus (ops); materialized views (prematuro).

## Decision: Errores en `pkg/httperr`

**Rationale**: Formato único `{"error":"...","code":404}`. Mapear `ErrNotFound`, `ErrForbidden`, `ErrValidation` en un handler gin.

**Alternatives considered**: Problema RFC 7807 (mejor, más verboso para el frontend aún inexistente).

## Decision: Docker multi-stage + Compose para Postgres

**Rationale**: Binario estático, imagen distroless/alpine &lt; 20–30 MB de app. Compose: api + postgres + volume uploads. `context.Context` desde Gin hasta `WithContext` de GORM.

**Alternatives considered**: Ejecutar Go en el host contra Postgres local (válido en dev, no es el artefacto de entrega).

## Resolved clarifications (ningún NEEDS CLARIFICATION pendiente)

| Tema | Resolución |
|------|------------|
| Frontend | Fuera de este repo |
| Tenant | `organizacion_id` |
| HTTP | Gin + `/api/v1` |
| Archivos | Interfaz + disco |
| Estados | Catálogo listado arriba |
| Billing | Fuera de v1 |
