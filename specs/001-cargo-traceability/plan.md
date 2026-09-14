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

**Performance Goals**: p95 &lt; 300 ms en CRUD y consulta de línea de tiempo con 1k cargas por organización (escala PyME, no high-frequency)

**Constraints**: JWT stateless; bcrypt cost ≥ 12; trazabilidad append-only; cancelación vía `context.Context` hasta GORM; errores JSON `{error, code}`; sin integraciones aduaneras/portuarias; aislamiento por `organizacion_id`

**Scale/Scope**: v1 = 1–N PyMEs (multi-tenant lógico), ~4 roles, decenas de usuarios y miles de cargas por tenant; facturación SaaS fuera de alcance

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

La constitution en `.specify/memory/constitution.md` **sigue siendo plantilla** (principios no ratificados). No hay gates formales que bloquear. Se aplican como **gates de trabajo** los principios del spec y de la arquitectura pedida:

| Gate | Estado | Nota |
|------|--------|------|
| Un despliegue, no microservicios | PASS | Un `cmd/api` |
| Dominios en `internal/` sin importar DB en handlers | PASS | Repositorios + interfaces |
| Contraseñas no en texto plano | PASS | bcrypt ≥ 12 |
| Autorización por rol y recurso (cliente = sus cargas) | PASS | JWT claims + middleware |
| Sin DIAN / terminal / naviera | PASS | Explicitamente fuera de contratos |
| Errores HTTP uniformes | PASS | `pkg/httperr` |
| Trazabilidad inmutable | PASS | Solo INSERT de eventos |
| Testeable (DI desde `main`) | PASS | Inyección de repos/servicios |

**Post-Phase 1**: sin violaciones nuevas. El patrón repositorio está justificado en Complexity Tracking (pedido explícito y constitution vacía).

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

Cada módulo de dominio: `handler.go`, `service.go`, `repository.go` (interfaz) + `postgres.go` (GORM). Los handlers **no** importan GORM.

**Structure Decision**: Backend-only en la raíz del repo (no `frontend/`). Layout estándar Go + módulos por bounded context alineados al spec (`usuarios`, `cargas`, `trazabilidad`, `documentos`, más `auth`, `clientes`, `incidencias`, `metricas`, `organizacion`). Un contenedor, una API `/api/v1`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Repository pattern (interfaces por módulo) | Constitution plantilla; el spec exige mocking y cambio de persistencia | Acceso GORM en handlers acopla HTTP a SQL y bloquea tests unitarios |
| Módulo `trazabilidad` separado de `cargas` | Auditoría append-only y consultas de línea de tiempo | Meter eventos solo en `cargas` mezcla mutación de estado con historial inmutable |
