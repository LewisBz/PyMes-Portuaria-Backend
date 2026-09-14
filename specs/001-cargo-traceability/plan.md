# Implementation Plan: API de gestión y trazabilidad de cargas (PyMEs portuarias)

**Branch**: `001-cargo-traceability` | **Date**: 2026-09-13 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-cargo-traceability/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## Summary

Backend REST SaaS para que una PyME logística portuaria gestione usuarios, clientes, cargas, estados, incidencias, documentos y métricas, con trazabilidad inmutable (quién, cuándo, estado anterior/nuevo). Alcance **solo operativa interna**: sin DIAN, terminales ni navieras.

Enfoque técnico: **monolito modular en Go** (Standard Go Project Layout), un despliegue, módulos por dominio en `internal/`, PostgreSQL + GORM + migraciones SQL (`golang-migrate`), JWT + bcrypt, repositorios por interfaz, Docker multi-stage. Este repo no incluye frontend.

Autores (contexto académico): Luis Manjarres, Hector Leon, Javier Solano, David Viaña — Universidad Libre Seccional Barranquilla (2026).

## Technical Context

**Language/Version**: Go 1.23+

**Primary Dependencies**: Gin (HTTP), GORM + driver PostgreSQL, golang-migrate, golang-jwt/jwt/v5, golang.org/x/crypto/bcrypt, godotenv (dev)

**Storage**: PostgreSQL 16; archivos de documentos en volumen local vía interfaz `ObjectStore` (evolución a S3/MinIO sin cambiar dominio)

**Testing**: `go test`, `net/http/httptest`, testify; tests unitarios con mocks de repositorio; tests de contrato sobre la API

**Target Platform**: Linux en contenedor (Docker); desarrollo en Windows con Docker Compose

**Project Type**: web-service (API REST JSON); monolito modular, un binario `cmd/api`

**Performance Goals**: p95 &lt; 300 ms en alta/listado/detalle de cargas y consulta de línea de tiempo con 1k cargas por organización (escala PyME, no high-frequency)

**Constraints**: JWT stateless con claims `user_id`, `organizacion_id`, `rol`; bcrypt cost ≥ 12; trazabilidad append-only; cancelación vía `context.Context` hasta GORM; errores JSON `{error, code}`; sin integraciones aduaneras/portuarias; aislamiento por `organizacion_id`

**Scale/Scope**: v1 = multi-tenant lógico con **una sola organización** vía seed/despliegue; ~4 roles; miles de cargas; facturación SaaS fuera de alcance

## Constitution Check

*GATE: Must pass against `.specify/memory/constitution.md` v1.0.0. Re-check after design.*

Evaluación formal de principios I–V de Constitution v1.0.0:

| Principle | Gate | Estado | Evidence |
|-----------|------|--------|----------|
| I. Modular Monolith First | Un `cmd/api`; dominios en `internal/`; handlers MUST NOT import GORM | PASS | Árbol de este plan; repositorios por interfaz (mandato I) |
| II. HTTP API Contract First | REST JSON `/api/v1` alineado a `contracts/openapi.yaml`; no CLI por paquete | PASS | OpenAPI de la feature; Compose/`go test` son entrypoints de operador |
| III. Contract Tests Before Implementation | `tests/contract/` en rojo antes de handlers por historia HTTP | PASS | tasks.md T017–T018, T031, T040, T047, T053, T060 |
| IV. Integration Tests for Persistence and AuthZ | Happy path, rol prohibido, y cross-cliente (404) | PASS | `tests/integration/` (T030, T039, T046, T052, T059, T063) |
| V. Simplicity, Tenant Isolation, Observability | Claims y filas con `organizacion_id`; eventos append-only; `{"error","code"}`; `context.Context` hasta GORM; sin DIAN | PASS | FR-001/013/014/015; T007–T011; T041; T067; T068 |

Security (JWT `user_id` + `organizacion_id` + `rol`, bcrypt ≥ 12, `organizaciones.activa` en middleware, ObjectStore): FR-001, T010, T025.

**Post-Phase 1**: PASS. El patrón repositorio no es violación: lo exige el principio I.

## Project Structure

### Documentation (this feature)

```text
specs/001-cargo-traceability/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/api/main.go                 # wiring: DB, migrate check, gin, DI
internal/
├── config/                     # env
├── platform/
│   ├── db/                     # *gorm.DB, context
│   ├── storage/                # ObjectStore (local)
│   └── httpx/                  # middleware JWT, RequireRole, errors
├── organizacion/
├── usuarios/
├── auth/
├── clientes/
├── cargas/
├── incidencias/
├── trazabilidad/
├── documentos/
└── metricas/
pkg/
└── httperr/                    # {"error","code"}
migrations/                     # NNN_name.up.sql / .down.sql
deploy/
├── Dockerfile
└── docker-compose.yml          # api + postgres
tests/
├── contract/                   # API httptest vs OpenAPI
├── integration/                # postgres testcontainer o compose
└── unit/                       # o tests junto a cada paquete *_test.go
```

Cada módulo de dominio: `model.go`, `handler.go`, `service.go`, `repository.go` (interfaz) + `postgres.go` (GORM). Los handlers **no** importan GORM.

**Structure Decision**: Backend-only en la raíz del repo (no `frontend/`). Layout estándar Go + módulos por bounded context alineados al spec (`usuarios`, `cargas`, `trazabilidad`, `documentos`, más `auth`, `clientes`, `incidencias`, `metricas`, `organizacion`). Un contenedor, una API `/api/v1`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

Ninguna violación. El patrón repositorio y el módulo `trazabilidad` separado son mandatorios (principios I y V), no excepciones.
