# PyMEs Portuaria — Backend

API REST de **gestión y trazabilidad de cargas** para PyMEs logísticas portuarias (Universidad Libre Seccional Barranquilla, 2026). Alcance v1: operativa interna de la PyME. **No** hay clientes hacia DIAN, terminales ni navieras.

Especificación ejecutable: [`specs/001-cargo-traceability/`](specs/001-cargo-traceability/). Guía Spec Kit: [`docs/spec-kit-guide.md`](docs/spec-kit-guide.md).

## Requisitos

- Go 1.23+
- Docker Desktop (PostgreSQL 16 + API)
- `gofmt` (incluido en Go). Opcional: `goimports`, [`golangci-lint`](https://golangci-lint.run/) con [`.golangci.yml`](.golangci.yml) (`govet`, `errcheck`, `staticcheck`)

```powershell
gofmt -w .
golangci-lint run
```

## Variables de entorno

Copia [`.env.example`](.env.example). `JWT_SECRET` ≥ 32 bytes. `BCRYPT_COST` ≥ 12. Semilla de desarrollo (migración `000003_seed_admin`):

| Campo | Valor |
|-------|--------|
| Organización | Demo PyME |
| Email | `admin@demo.local` |
| Password | `DemoAdmin12!` |

No uses esta contraseña en producción.

## Levantar con Compose

```powershell
docker compose -f deploy/docker-compose.yml up --build -d
curl -s http://localhost:8080/healthz
```

Base URL de la API: `http://localhost:8080/api/v1`.

Contrato OpenAPI: [`specs/001-cargo-traceability/contracts/openapi.yaml`](specs/001-cargo-traceability/contracts/openapi.yaml). Escenarios manuales: [`specs/001-cargo-traceability/quickstart.md`](specs/001-cargo-traceability/quickstart.md).

## Tests

```powershell
go test ./...
```

Los tests de `tests/integration` se saltan si no hay `DATABASE_URL` o `TEST_DATABASE_URL`.

## Arquitectura

Monolito modular: un binario `cmd/api`, dominios en `internal/`, handlers sin GORM, JWT claims `user_id`, `organizacion_id`, `rol`.
