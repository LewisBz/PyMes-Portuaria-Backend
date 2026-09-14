CREATE TYPE usuario_rol AS ENUM ('administrador', 'operador', 'gerente', 'cliente');

CREATE TABLE usuarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizacion_id UUID NOT NULL REFERENCES organizaciones(id) ON DELETE RESTRICT,
    email VARCHAR(320) NOT NULL,
    password_hash TEXT NOT NULL,
    nombre VARCHAR(200) NOT NULL,
    rol usuario_rol NOT NULL,
    cliente_id UUID NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organizacion_id, email)
);

CREATE INDEX idx_usuarios_org_email ON usuarios (organizacion_id, email);
