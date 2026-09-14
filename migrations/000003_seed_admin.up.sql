-- Seed tenant Demo PyME + administrador.
-- Password (dev only): DemoAdmin12!
-- Hash generated with golang.org/x/crypto/bcrypt cost 12.

INSERT INTO organizaciones (id, nombre, activa)
VALUES ('00000000-0000-4000-8000-000000000001', 'Demo PyME', TRUE);

INSERT INTO usuarios (id, organizacion_id, email, password_hash, nombre, rol, activo)
VALUES (
    '00000000-0000-4000-8000-000000000002',
    '00000000-0000-4000-8000-000000000001',
    'admin@demo.local',
    '$2a$12$wKwU1uTjHbtYsbk99M4CxemOqz05SeFOI02uAaCdCz1Y5mv7W.HiK',
    'Administrador Demo',
    'administrador',
    TRUE
);
