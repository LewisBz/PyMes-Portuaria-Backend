CREATE TYPE incidencia_estado AS ENUM ('abierta', 'cerrada');

CREATE TABLE incidencias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizacion_id UUID NOT NULL REFERENCES organizaciones(id) ON DELETE RESTRICT,
    carga_id UUID NOT NULL REFERENCES cargas(id) ON DELETE RESTRICT,
    titulo VARCHAR(200) NOT NULL,
    descripcion TEXT NOT NULL,
    tipo VARCHAR(40) NOT NULL DEFAULT 'otro',
    estado incidencia_estado NOT NULL DEFAULT 'abierta',
    abierta_por UUID NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    cerrada_por UUID REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ
);

CREATE INDEX idx_incidencias_carga ON incidencias (carga_id);
