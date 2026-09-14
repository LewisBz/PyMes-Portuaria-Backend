# Data model: 001-cargo-traceability

Todas las tablas de negocio (salvo `organizaciones`) incluyen `organizacion_id`. Soft-delete solo en usuarios/clientes (`activo`). Trazabilidad **sin** updates.

## Entity: Organizacion

Representa la PyME suscriptora (tenant).

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | gen_random_uuid() |
| nombre | varchar(200) | required |
| activa | boolean | default true; si false, login bloqueado |
| created_at, updated_at | timestamptz | |

**Relationships**: 1:N Usuario, Cliente, Carga.

## Entity: Usuario

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | |
| organizacion_id | UUID FK | required, index |
| email | varchar(320) | unique por organización; login |
| password_hash | text | bcrypt; nunca exponer en JSON |
| nombre | varchar(200) | required |
| rol | enum | `administrador`, `operador`, `gerente`, `cliente` |
| cliente_id | UUID FK nullable | required si rol = `cliente`; el portal de ese cliente |
| activo | boolean | default true |
| created_at, updated_at | timestamptz | |

**Validation**: email formato RFC; password en claro solo en input (min 8), nunca persistido; un único `cliente_id` por usuario cliente.

**State**: `activo` true/false (desactivar ≠ borrar, para FK de trazabilidad).

## Entity: Cliente

Cliente comercial de la PyME (no confundir con rol).

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | |
| organizacion_id | UUID FK | |
| razon_social | varchar(200) | required |
| nit | varchar(20) | opcional, unique por org si presente |
| email | varchar(320) | opcional |
| activo | boolean | default true |
| created_at, updated_at | timestamptz | |

**Relationships**: 1:N Carga; 0:N Usuario (rol cliente).

## Entity: Carga

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | |
| organizacion_id | UUID FK | |
| cliente_id | UUID FK | required, mismo tenant |
| referencia | varchar(80) | unique por org; identificador interno |
| descripcion | text | opcional |
| origen | varchar(200) | opcional |
| destino | varchar(200) | opcional |
| estado | enum | ver catálogo |
| fecha_comprometida | date | nullable; base de “retraso” |
| created_at, updated_at | timestamptz | |

**Estado catálogo**: `registrada` (alta), `en_transito`, `en_puerto`, `entregada`, `cancelada`.

**Derived**: `retrasada` = `fecha_comprometida < today` AND estado ∉ (`entregada`, `cancelada`). No se persiste; se expone en API y métricas.

**State transitions**: desde no terminal hacia otro distinto; mismo estado → 409. Terminal (`entregada`, `cancelada`) no cambia salvo administrador (v1: **no** se reabre).

## Entity: Incidencia

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | |
| organizacion_id | UUID FK | |
| carga_id | UUID FK | required |
| titulo | varchar(200) | required |
| descripcion | text | required |
| tipo | varchar(40) | p. ej. `demora`, `dano`, `documental`, `otro` |
| estado | enum | `abierta`, `cerrada` |
| abierta_por | UUID FK usuarios | |
| cerrada_por | UUID FK nullable | |
| created_at, closed_at | timestamptz | closed_at null si abierta |

**Transitions**: abierta → cerrada (no reabrir en v1).

## Entity: Documento

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | |
| organizacion_id | UUID FK | |
| carga_id | UUID FK | |
| nombre_original | varchar(255) | |
| content_type | varchar(100) | `application/pdf`, `image/jpeg`, `image/png` |
| tamano_bytes | int | 1..10_485_760 |
| storage_key | varchar(500) | path interno, no URL pública cruda |
| visible_cliente | boolean | default false |
| subido_por | UUID FK | |
| created_at | timestamptz | |

## Entity: EventoTrazabilidad

Inmutable. Sin `updated_at`. Sin DELETE en aplicación.

| Field | Type | Rules |
|-------|------|--------|
| id | UUID PK | |
| organizacion_id | UUID FK | |
| carga_id | UUID FK | index (carga_id, occurred_at) |
| tipo | enum | `estado_cambiado`, `incidencia_abierta`, `incidencia_cerrada`, `documento_adjunto`, `carga_creada` |
| actor_id | UUID FK usuarios | |
| occurred_at | timestamptz | default now() |
| estado_anterior | varchar(40) | null si no aplica |
| estado_nuevo | varchar(40) | null si no aplica |
| payload | jsonb | ids relacionados, notas cortas |

**Rule**: todo cambio de estado de carga o de incidencia **MUST** insertar un evento en la misma transacción.

## Relationships (overview)

```text
Organizacion
  ├── Usuario (rol; si cliente → Cliente)
  ├── Cliente
  │     └── Carga
  │           ├── Incidencia
  │           ├── Documento
  │           └── EventoTrazabilidad
```

## Integrity

- FK con `ON DELETE RESTRICT` en cargas/eventos (no borrar historia).
- Índices: `(organizacion_id, email)` usuarios; `(organizacion_id, referencia)` cargas; `(carga_id, occurred_at)` eventos.
- Toda query de listado filtra `organizacion_id` del JWT; si rol `cliente`, además `cargas.cliente_id = usuario.cliente_id`.
