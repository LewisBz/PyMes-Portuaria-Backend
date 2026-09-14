---
description: "Task list for 001-cargo-traceability"
---

#    Tasks: API de gestión y trazabilidad de cargas (PyMEs portuarias)

**Input**: Design documents from `/specs/001-cargo-traceability/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: Incluidas. El spec exige pruebas de contrato por rol (SC-004) y el plan define `tests/contract/` + `tests/integration/`.

**Organization**: Tareas agrupadas por historia de usuario para implementación y prueba independientes.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Puede correr en paralelo (archivos distintos, sin depender de tareas incompletas)
- **[Story]**: US1–US6 según spec.md
- Cada descripción incluye ruta de archivo exacta



## Path Conventions

Monolito Go en la raíz del repo (`cmd/api`, `internal/`, `pkg/`, `migrations/`, `deploy/`, `tests/`) según plan.md.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Inicializar el módulo Go y el esqueleto de directorios del monolito modular

- [ ] T001 Create empty package directories `cmd/api/`, `internal/config/`, `internal/platform/db/`, `internal/platform/storage/`, `internal/platform/httpx/`, `internal/organizacion/`, `internal/usuarios/`, `internal/auth/`, `internal/clientes/`, `internal/cargas/`, `internal/incidencias/`, `internal/trazabilidad/`, `internal/documentos/`, `internal/metricas/`, `pkg/httperr/`, `migrations/`, `deploy/`, `tests/contract/`, `tests/integration/`
- [ ] T002 Initialize `go.mod` (Go 1.23+) and require `github.com/gin-gonic/gin`, `gorm.io/gorm`, `gorm.io/driver/postgres`, `github.com/golang-migrate/migrate/v4`, `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto`, `github.com/joho/godotenv`, `github.com/stretchr/testify`
- [ ] T003 [P] Add `.golangci.yml` with govet, errcheck, staticcheck and document `gofmt`/`goimports` in `README.md`
- [ ] T004 [P] Add `.env.example` with `DATABASE_URL`, `JWT_SECRET`, `JWT_EXPIRY`, `PORT`, `UPLOAD_DIR`, `BCRYPT_COST=12`
- [ ] T005 [P] Add root `.gitignore` entries for `.env`, `bin/`, `data/uploads/`, `*.exe`, `graphify-out/`

---



## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Config, DB, errores HTTP, JWT/roles, migraciones, Compose y `main` — bloquea todas las historias

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 Implement env loading and typed config in `internal/config/config.go`
- [ ] T007 Implement uniform JSON errors `{"error","code"}` in `pkg/httperr/httperr.go`
- [ ] T008 Open PostgreSQL with GORM and pass `context.Context` via `WithContext` in `internal/platform/db/db.go`
- [ ] T009 [P] Map `httperr` to Gin JSON responses in `internal/platform/httpx/errors.go`
- [ ] T010 [P] Implement JWT parse/validate middleware (claims `user_id`, `organizacion_id`, `rol`) and reject the request with 401 when `organizaciones.activa` is false in `internal/platform/httpx/jwt.go`
- [ ] T011 [P] Implement `RequireRole` middleware in `internal/platform/httpx/roles.go`
- [ ] T012 Run `golang-migrate` on startup in `internal/platform/db/migrate.go` reading `migrations/`
- [ ] T013 Add `migrations/000001_organizaciones.up.sql` and `migrations/000001_organizaciones.down.sql` (table `organizaciones` per data-model.md)
- [ ] T014 Wire Gin engine, `/healthz`, `/api/v1` group, and DI stubs in `cmd/api/main.go`
- [ ] T015 Add PostgreSQL 16 service and `DATABASE_URL` in `deploy/docker-compose.yml`
- [ ] T016 Add multi-stage static binary image in `deploy/Dockerfile`

**Checkpoint**: Foundation ready — `docker compose -f deploy/docker-compose.yml up` arranca Postgres; `/healthz` responde; user stories pueden empezar

---



## Phase 3: User Story 1 - Autenticación y administración de usuarios (Priority: P1) 🎯 MVP

**Goal**: Login JWT, cuatro roles, CRUD/desactivar usuarios por administrador; cliente no lista usuarios

**Independent Test**: Semilla admin → login → crear operador y cliente → token de cliente recibe 403 en `GET /api/v1/usuarios`; credenciales inválidas 401 sin filtrar si el correo existe. El JWT incluye `user_id`, `organizacion_id` y `rol`. (El caso “cliente no crea cargas” está en US2.)

### Tests for User Story 1

> Write these tests FIRST, ensure they FAIL before implementation

- [ ] T017 [P] [US1] Contract tests for `POST /api/v1/auth/login` and `GET /api/v1/auth/me` asserting JWT claims `user_id`, `organizacion_id`, `rol` in `tests/contract/auth_test.go` against `specs/001-cargo-traceability/contracts/openapi.yaml`
- [ ] T018 [P] [US1] Contract tests for `GET/POST /api/v1/usuarios` and `PATCH /api/v1/usuarios/{id}` (admin vs cliente 403) in `tests/contract/usuarios_test.go`



### Implementation for User Story 1

- [ ] T019 [US1] Add `migrations/000002_usuarios.up.sql` and `migrations/000002_usuarios.down.sql` (`usuarios`, unique `(organizacion_id, email)`, rol enum, `cliente_id` nullable)
- [ ] T020 [P] [US1] Define `Usuario` and `Rol` types in `internal/usuarios/model.go`
- [ ] T021 [P] [US1] Define `Organizacion` type and `Repository` interface in `internal/organizacion/model.go` and `internal/organizacion/repository.go`
- [ ] T022 [US1] Implement GORM `Organizacion` repository in `internal/organizacion/postgres.go` (handlers must not import GORM)
- [ ] T023 [US1] Implement `Repository` interface and GORM adapter in `internal/usuarios/repository.go` and `internal/usuarios/postgres.go`
- [ ] T024 [US1] Hash passwords with bcrypt cost ≥ 12 and implement create/list/patch/deactivate in `internal/usuarios/service.go`
- [ ] T025 [US1] Issue JWT (`user_id`, `organizacion_id`, `rol`) and reject inactive user **or inactive organization** (`organizaciones.activa = false`) in `internal/auth/service.go`
- [ ] T026 [US1] Implement `POST /auth/login` and `GET /auth/me` in `internal/auth/handler.go`
- [ ] T027 [US1] Implement usuario routes with `RequireRole("administrador")` in `internal/usuarios/handler.go`
- [ ] T028 [US1] Register auth and usuarios routes and inject repos/services in `cmd/api/main.go`
- [ ] T029 [US1] Seed org + admin (`admin@demo.local`) in `migrations/000003_seed_admin.up.sql` and `migrations/000003_seed_admin.down.sql`
- [ ] T030 [US1] Add login/role integration coverage in `tests/integration/auth_roles_test.go`

**Checkpoint**: US1 usable and independently testable (MVP auth)

---



## Phase 4: User Story 2 - Clientes de la PyME y registro de cargas (Priority: P1)

**Goal**: Operador registra clientes y cargas (`estado=registrada`); listados filtrados por tenant; cliente solo ve sus cargas

**Independent Test**: Operador crea cliente + carga y la lista; usuario cliente **no puede crear cargas** (`POST /api/v1/cargas` → 403); `GET /api/v1/cargas/{id}` de carga de otro cliente → **404** (no 403); carga sin `cliente_id` válido → 400

### Tests for User Story 2

- [ ] T031 [P] [US2] Contract tests for `GET/POST /api/v1/clientes` and `GET/POST /api/v1/cargas` in `tests/contract/clientes_cargas_test.go`



### Implementation for User Story 2

- [ ] T032 [US2] Add `migrations/000004_create_clientes_cargas.up.sql` and `migrations/000004_create_clientes_cargas.down.sql` containing **only** table `clientes` (FK tenant) in this task
- [ ] T033 [P] [US2] Implement cliente model, repository interface, GORM adapter, and service in `internal/clientes/model.go`, `internal/clientes/repository.go`, `internal/clientes/postgres.go`, `internal/clientes/service.go`
- [ ] T034 [US2] Implement cliente HTTP handlers in `internal/clientes/handler.go`
- [ ] T035 [US2] Implement carga model, repository, GORM adapter, and service (create/list/get + tenant and `cliente_id` scope; **no** persistir eventos de trazabilidad) in `internal/cargas/model.go`, `internal/cargas/repository.go`, `internal/cargas/postgres.go`, `internal/cargas/service.go`
- [ ] T036 [US2] Add table `cargas` only (enum SQL inmutable, unique `(organizacion_id, referencia)`, `fecha_comprometida`, `fecha_entrega` nullable) to `migrations/000004_create_clientes_cargas.up.sql` and `.down.sql`. Do **not** create `eventos_trazabilidad` (that is exclusively T041 / `000005_create_eventos_trazabilidad`)
- [ ] T037 [US2] Implement carga HTTP handlers in `internal/cargas/handler.go`
- [ ] T038 [US2] Register clientes and cargas routes in `cmd/api/main.go`
- [ ] T039 [US2] Add scope tests in `tests/integration/cargas_scope_test.go`: other-cliente `GET /cargas/{id}` returns **404 Not Found** (never 403); cliente `POST /cargas` returns 403

**Checkpoint**: US1 + US2 independently testable; operador puede dar de alta operación

---



## Phase 5: User Story 3 - Estados de carga y línea de tiempo (Priority: P1)

**Goal**: Cambio de estado en transacción con evento inmutable; GET línea de tiempo; cliente no cambia estados

**Independent Test**: Dos cambios de estado → timeline ordenada con actor y estados anterior/nuevo; **gerente GET trazabilidad → 200**; mismo estado → 409; cliente POST estados → 403

### Tests for User Story 3

- [ ] T040 [P] [US3] Contract tests for `POST /api/v1/cargas/{id}/estados` and `GET /api/v1/cargas/{id}/trazabilidad` (include **gerente → 200 OK**) in `tests/contract/trazabilidad_test.go`



### Implementation for User Story 3

- [ ] T041 [US3] Add `migrations/000005_create_eventos_trazabilidad.up.sql` and `migrations/000005_create_eventos_trazabilidad.down.sql` as the **sole** DDL owner of table `eventos_trazabilidad` (append-only, index `(carga_id, occurred_at)`). Do not create this table in any other migration.
- [ ] T042 [US3] Implement evento model, insert-only postgres adapter, chronological list, and `Append` (including `carga_creada` from carga service after this table exists) in `internal/trazabilidad/model.go`, `internal/trazabilidad/postgres.go`, `internal/trazabilidad/service.go`, `internal/trazabilidad/repository.go`
- [ ] T043 [US3] Implement `CambiarEstado` (reject same/terminal states; set `fecha_entrega` when transitioning to `entregada`) in the same DB transaction as event insert in `internal/cargas/service.go`
- [ ] T044 [US3] Add `POST /cargas/{id}/estados` in `internal/cargas/handler.go` and `GET /cargas/{id}/trazabilidad` in `internal/trazabilidad/handler.go` with role/resource checks
- [ ] T045 [US3] Register trazabilidad routes in `cmd/api/main.go`
- [ ] T046 [US3] Assert one event per state change in `tests/integration/trazabilidad_tx_test.go`

**Checkpoint**: Problema central de trazabilidad cubierto a nivel API

---



## Phase 6: User Story 4 - Incidencias operativas (Priority: P2)

**Goal**: Abrir/cerrar incidencias ligadas a carga; eventos `incidencia_abierta` / `incidencia_cerrada`; cliente ve, no cierra

**Independent Test**: Abrir incidencia → listar en carga + evento; cerrar → `cerrada` + evento; cliente POST → 403

### Tests for User Story 4

- [ ] T047 [P] [US4] Contract tests for `GET/POST /api/v1/cargas/{id}/incidencias` and `POST /api/v1/incidencias/{id}/cerrar` in `tests/contract/incidencias_test.go`



### Implementation for User Story 4

- [ ] T048 [US4] Add `migrations/000006_incidencias.up.sql` and `migrations/000006_incidencias.down.sql`
- [ ] T049 [US4] Implement model, repository, GORM adapter, and service (open/close, no reopen) in `internal/incidencias/model.go`, `internal/incidencias/repository.go`, `internal/incidencias/postgres.go`, `internal/incidencias/service.go`
- [ ] T050 [US4] Emit trazabilidad events inside incidencia transactions in `internal/incidencias/service.go`
- [ ] T051 [US4] Implement handlers in `internal/incidencias/handler.go` and register routes in `cmd/api/main.go`
- [ ] T052 [US4] Add integration coverage in `tests/integration/incidencias_test.go`

**Checkpoint**: US4 independently testable on top of cargas + trazabilidad

---



## Phase 7: User Story 5 - Documentos adjuntos (Priority: P2)

**Goal**: Upload PDF/JPEG/PNG ≤ 10 MB; `visible_cliente`; cliente solo lista/descarga autorizados

**Independent Test**: Upload visible → cliente download 200; interno → cliente 403; `.exe` o >10 MB → 400

### Tests for User Story 5

- [ ] T053 [P] [US5] Contract tests for `GET/POST /api/v1/cargas/{id}/documentos` and `GET /api/v1/documentos/{id}/download` in `tests/contract/documentos_test.go`



### Implementation for User Story 5

- [ ] T054 [US5] Implement `ObjectStore` (`Put`, `Get`, `Delete`) with local disk in `internal/platform/storage/store.go` and `internal/platform/storage/local.go`
- [ ] T055 [US5] Add `migrations/000007_documentos.up.sql` and `migrations/000007_documentos.down.sql`
- [ ] T056 [US5] Implement model, repository, service (MIME/size validation, `visible_cliente`) in `internal/documentos/model.go`, `internal/documentos/repository.go`, `internal/documentos/postgres.go`, `internal/documentos/service.go`
- [ ] T057 [US5] Implement multipart upload, list, and download handlers in `internal/documentos/handler.go` and register them in `cmd/api/main.go`
- [ ] T058 [US5] Append `documento_adjunto` events in `internal/documentos/service.go`
- [ ] T059 [US5] Add upload/download authorization tests in `tests/integration/documentos_test.go`

**Checkpoint**: Portal cliente puede ver documentos autorizados sin sistemas aduaneros

---



## Phase 8: User Story 6 - Tablero de métricas para gerencia (Priority: P3)

**Goal**: `GET /metricas` para administrador y gerente; 403 para cliente

**Independent Test**: Con cargas e incidencias de prueba, agregados coherentes; cliente 403

### Tests for User Story 6

- [ ] T060 [P] [US6] Contract tests for `GET /api/v1/metricas` (gerente 200, cliente 403) in `tests/contract/metricas_test.go`



### Implementation for User Story 6

- [ ] T061 [US6] Implement SQL aggregates in `internal/metricas/service.go` and `internal/metricas/postgres.go`: `cargas_activas`, `cargas_retrasadas` (derived `fecha_comprometida < NOW() AND estado != 'entregada'`), `incidencias_abiertas`, and **tiempo promedio total de entrega** in hours (`AVG(fecha_entrega - created_at)` only for `estado = 'entregada'`). Do not compute inter-state intervals in v1.
- [ ] T062 [US6] Implement `GET /metricas` with `RequireRole("administrador","gerente")` in `internal/metricas/handler.go` and register in `cmd/api/main.go`
- [ ] T063 [US6] Add aggregation tests in `tests/integration/metricas_test.go`

**Checkpoint**: Todas las historias del spec son demostrables por API

---



## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Entrega, seguridad y validación del quickstart

- [ ] T064 [P] Document setup, seed admin, Compose, API base URL, and Spec Kit feature path `specs/001-cargo-traceability/` (link `docs/spec-kit-guide.md`) in `README.md`
- [ ] T065 Execute scenarios in `specs/001-cargo-traceability/quickstart.md` against local Compose and record gaps
- [ ] T066 Confirm bcrypt cost ≥ 12 and JWT expiry from `internal/config/config.go` / `internal/auth/service.go`
- [ ] T067 Grep the repo to ensure no DIAN, terminal, or naviera clients exist under `internal/`
- [ ] T068 Ensure Gin request `context.Context` is passed into GORM in `internal/platform/db/db.go` and all `postgres.go` adapters

---



## Dependencies & Execution Order



### Phase Dependencies

- **Setup (Phase 1)**: Sin dependencias
- **Foundational (Phase 2)**: Depende de Setup; **bloquea** todas las historias
- **US1 (Phase 3)**: Depende de Phase 2 — MVP
- **US2 (Phase 4)**: Depende de US1 (roles operador/cliente)
- **US3 (Phase 5)**: Depende de US2 (cargas existentes)
- **US4 (Phase 6)**: Depende de US2 + US3 (carga + Append de eventos)
- **US5 (Phase 7)**: Depende de US2 (carga) + US3 (eventos)
- **US6 (Phase 8)**: Depende de US2 y US4 para datos agregables
- **Polish (Phase 9)**: Depende de las historias que se quieran entregar



### User Story Dependencies

- **US1 (P1)**: Tras Phase 2; no depende de otras historias
- **US2 (P1)**: Tras US1 (autenticación)
- **US3 (P1)**: Tras US2
- **US4 (P2)**: Tras US3 (reutiliza trazabilidad)
- **US5 (P2)**: Tras US2; eventos si US3 listo
- **US6 (P3)**: Tras US2; mejor con US4 para incidencias abiertas



### Within Each User Story

- Tests de contrato primero y en rojo
- Model/migración → repositorio → servicio → handler → `main.go` → integración
- Completar checkpoint antes de subir de prioridad



### Parallel Opportunities

- T003, T004, T005 en paralelo (Phase 1)
- T009, T010, T011 en paralelo (Phase 2)
- T017 || T018 (contratos US1)
- T020 || T021 (modelos US1)
- T033 modelos clientes vs inicio de contratos US2 (T031) en paralelo al arrancar la fase
- T047–T053–T060: contratos de US4/US5/US6 en paralelo **después** de sus fases de datos, o en rojo al inicio de cada fase
- Un desarrollador no debe editar el mismo `cmd/api/main.go` en paralelo sin integrar

---



## Parallel Example: User Story 1

```bash
# Contratos en paralelo (deben fallar):
Task: "Contract tests POST /auth/login and GET /auth/me in tests/contract/auth_test.go"
Task: "Contract tests usuarios admin vs cliente in tests/contract/usuarios_test.go"

# Tipos en paralelo:
Task: "Usuario/Rol in internal/usuarios/model.go"
Task: "Organizacion types in internal/organizacion/model.go"
```



## Parallel Example: User Story 5 vs 6 (after US2/US3)

```bash
# Solo si el equipo se parte por archivos distintos:
Task: "ObjectStore in internal/platform/storage/local.go"
Task: "Metricas SQL in internal/metricas/postgres.go"
```

---



## Implementation Strategy



### MVP First (User Story 1 Only)

1. Completar Phase 1: Setup
2. Completar Phase 2: Foundational
3. Completar Phase 3: US1
4. **STOP and VALIDATE**: login, roles, 403 de cliente en `/usuarios`
5. Demo académica de seguridad



### Incremental Delivery

1. Setup + Foundational
2. US1 → demo auth (MVP)
3. US2 → altas operativas
4. US3 → trazabilidad (valor central del grado)
5. US4 → incidencias
6. US5 → documentos
7. US6 → tablero gerencial
8. Polish + `quickstart.md`



### Parallel Team Strategy

1. Equipo junto en Phase 1–2
2. Tras foundation: prioridad US1 (un owner de `cmd/api/main.go`)
3. US2–US3 en serie (mismo agregado Carga)
4. US4 y US5 en paralelo solo con owners distintos (`internal/incidencias/` vs `internal/documentos/` + `storage/`)
5. US6 al final sobre datos reales de prueba

---



## Notes

- [P] = archivos distintos y sin dependencia de tareas incompletas
- Handlers **nunca** importan `gorm.io/gorm`
- No hay tareas de frontend ni de facturación SaaS ni de DIAN
- Commit por tarea o grupo lógico
- Validar cada checkpoint con el Independent Test de la historia

