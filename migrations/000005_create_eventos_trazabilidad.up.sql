CREATE TYPE evento_tipo AS ENUM (
    'estado_cambiado',
    'incidencia_abierta',
    'incidencia_cerrada',
    'documento_adjunto',
    'carga_creada'
);

CREATE TABLE eventos_trazabilidad (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizacion_id UUID NOT NULL REFERENCES organizaciones(id) ON DELETE RESTRICT,
    carga_id UUID NOT NULL REFERENCES cargas(id) ON DELETE RESTRICT,
    tipo evento_tipo NOT NULL,
    actor_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    estado_anterior VARCHAR(40),
    estado_nuevo VARCHAR(40),
    payload JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX idx_eventos_carga_occurred ON eventos_trazabilidad (carga_id, occurred_at);
