CREATE TABLE documentos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizacion_id UUID NOT NULL REFERENCES organizaciones(id) ON DELETE RESTRICT,
    carga_id UUID NOT NULL REFERENCES cargas(id) ON DELETE RESTRICT,
    nombre_original VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    tamano_bytes INT NOT NULL CHECK (tamano_bytes >= 1 AND tamano_bytes <= 10485760),
    storage_key VARCHAR(500) NOT NULL,
    visible_cliente BOOLEAN NOT NULL DEFAULT FALSE,
    subido_por UUID NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documentos_carga ON documentos (carga_id);
