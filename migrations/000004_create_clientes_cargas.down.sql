DROP TABLE IF EXISTS cargas;
DROP TYPE IF EXISTS carga_estado;
ALTER TABLE usuarios DROP CONSTRAINT IF EXISTS usuarios_cliente_id_fkey;
DROP TABLE IF EXISTS clientes;
