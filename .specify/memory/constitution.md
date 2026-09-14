<!--
Sync Impact Report
- Version change: (unratified template) → 1.0.0
- Modified principles:
  - template principle 1 → I. Modular Monolith First
  - template principle 2 (CLI) → II. HTTP API Contract First
  - template principle 3 (TDD) → III. Contract Tests Before Implementation
  - template principle 4 → IV. Integration Tests for Persistence and AuthZ
  - template principle 5 → V. Simplicity, Tenant Isolation, and Honest Observability
- Added sections: Security and Scope Constraints; Review Gates
- Removed sections: none (template placeholders replaced)
- Follow-up TODOs: none
-->

# PyMEs Portuaria Backend Constitution

## Core Principles

### I. Modular Monolith First

This repository delivers **one** deployable API (`cmd/api`). Domain logic MUST live
in `internal/` packages by bounded context (usuarios, clientes, cargas,
trazabilidad, incidencias, documentos, metricas). HTTP handlers MUST NOT import
GORM or SQL drivers; they MUST depend on service and repository **interfaces**.
New microservices, extra binaries for the same domain, or a second persistence
style for the same aggregate REQUIRE a constitution amendment plus Complexity
Tracking in the feature plan.

Rationale: a PyME academic SaaS cannot absorb distributed-system cost. Package
boundaries keep the codebase testable without splitting deploys.

### II. HTTP API Contract First

User-visible behavior MUST be expressed as versioned REST JSON under `/api/v1`
and MUST match `specs/*/contracts/` (OpenAPI) for the active feature. Breaking
path, status, or error-shape changes MUST bump the API version or be rejected.
This project MUST NOT require a CLI per package. Operator and CI entrypoints
MAY exist (`go test`, Docker Compose) but they are not the product surface.

Rationale: the product is a backend for a web SaaS, not a library-plus-CLI
toolkit. Contracts keep frontend and tests aligned.

### III. Contract Tests Before Implementation

For each user story that exposes HTTP, contract tests under `tests/contract/`
MUST be written first and MUST fail against the unimplemented API before
handlers are filled in. Implementation MUST then make those tests pass without
weakening assertions to match accidental behavior.

Rationale: Spec Kit stories are independently demonstrable. Red-then-green on
the OpenAPI surface is the enforceable TDD gate for this repo (not a separate
human “approve the test file” ceremony).

### IV. Integration Tests for Persistence and AuthZ

Changes to schema, JWT/role rules, or tenant/cliente scoping MUST include
integration tests under `tests/integration/` (or package tests that hit Postgres
in Compose). Those tests MUST cover at least: happy path, forbidden role, and
cross-tenant or cross-cliente access. Contract tests do not replace this gate.

Rationale: authorization bugs and missing `organizacion_id` filters will not
show up in mocked unit tests.

### V. Simplicity, Tenant Isolation, and Honest Observability

Features MUST stay inside operativa interna de la PyME. Every business row
MUST carry `organizacion_id` and queries MUST filter it from JWT claims.
Traceability events MUST be append-only. Prefer the smallest design that
satisfies the spec (YAGNI). HTTP errors MUST use `{"error","code"}`. Request
`context.Context` MUST propagate to GORM so cancelled clients do not keep
running queries. Structured logs SHOULD include request id, org id, and user
id; they MUST NOT log passwords, JWT secrets, or document bytes.

Rationale: isolation and auditability are the academic and operational core;
complexity and leaked secrets are not.

## Security and Scope Constraints

Stack for this constitution: Go 1.23+, Gin, PostgreSQL 16, GORM adapters behind
interfaces, golang-migrate SQL files, JWT (claims MUST include user id,
organizacion id, and role), bcrypt cost MUST be >= 12. Passwords MUST never be
stored or returned in plaintext.

The API MUST NOT integrate with DIAN, terminales portuarios, navieras, or other
official customs/port systems. Subscription billing is out of scope until a
feature spec says otherwise.

Documents MUST be stored behind an `ObjectStore` interface; v1 MAY use local
disk. Allowed types and size limits are defined by the active spec (PDF, JPEG,
PNG, 10 MB unless the spec is amended).

## Review Gates

A change is ready to merge only if:

1. It complies with Core Principles I–V and Security and Scope Constraints.
2. `go test ./...` passes, including new contract/integration tests for the story.
3. Handlers still do not import GORM.
4. No new DIAN/terminal/naviera client appears under `internal/`.
5. Added complexity (extra service, extra datastore, extra deploy) is justified
   in the feature `plan.md` Complexity Tracking table.

Runtime agent guidance: Spec Kit artifacts under `specs/` and this constitution.
Do not treat unratified template comments in other files as law.

## Governance

This constitution supersedes informal chat decisions and example Spec Kit
template principles (library-first, CLI-per-library). Specs, plans, and tasks
MUST be adjusted to fit this document; this document MUST NOT be silently
ignored to fit a spec.

Amendments: edit `.specify/memory/constitution.md`, bump version, set
`Last Amended` to the amendment date, and record a Sync Impact Report comment.
MAJOR: remove or invert a MUST principle. MINOR: add a principle or a
normative constraint. PATCH: wording-only.

Ratification date is the first adoption of v1.0.0. Compliance review happens in
`/speckit-plan` Constitution Check, `/speckit-analyze`, and pull-request review.
Unresolved MUST violations block `/speckit-implement` for the affected story.

**Version**: 1.0.0 | **Ratified**: 2026-09-13 | **Last Amended**: 2026-09-13
