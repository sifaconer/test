-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clientes (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR NOT NULL
);
CREATE TABLE IF NOT EXISTS productos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR NOT NULL,
    precio DECIMAL NOT NULL
);
CREATE TABLE IF NOT EXISTS carrito_compra (
    id SERIAL PRIMARY KEY,
    cliente_id INT REFERENCES clientes(id),
    producto_id INT REFERENCES productos(id),
    cantidad INT NOT NULL,
    fecha_agregado TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS carrito_compra;
DROP TABLE IF EXISTS clientes;
DROP TABLE IF EXISTS productos;
-- +goose StatementEnd
