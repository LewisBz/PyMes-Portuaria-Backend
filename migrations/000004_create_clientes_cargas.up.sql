CREATE TABLE clientes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizacion_id UUID NOT NULL REFERENCES organizaciones(id) ON DELETE RESTRICT,
    razon_social VARCHAR(200) NOT NULL,
    nit VARCHAR(20),
    email VARCHAR(320),
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organizacion_id, nit)
);

CREATE INDEX idx_clientes_org ON clientes (organizacion_id);

ALTER TABLE usuarios
    ADD CONSTRAINT usuarios_cliente_id_fkey
    FOREIGN KEY (cliente_id) REFERENCES clientes(id) ON DELETE RESTRICT;

CREATE TYPE carga_estado AS ENUM ('registrada', 'en_transito', 'en_puerto', 'entregada', 'cancelada');

CREATE TABLE cargas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizacion_id UUID NOT NULL REFERENCES organizaciones(id) ON DELETE RESTRICT,
    cliente_id UUID NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    referencia VARCHAR(80) NOT NULL,
    descripcion TEXT,
    origen VARCHAR(200),
    destino VARCHAR(200),
    estado carga_estado NOT NULL DEFAULT 'registrada',
    fecha_comprometida DATE,
    fecha_entrega TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organizacion_id, referencia)
);

CREATE INDEX idx_cargas_org_ref ON cargas (organizacion_id, referencia);
CREATE INDEX idx_cargas_cliente ON cargas (cliente_id);
