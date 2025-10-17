-- name: ListProducts :many
SELECT id, name, price, created_at FROM products ORDER BY created_at DESC;

-- name: CreateProduct :one
INSERT INTO products (name, price)
VALUES ($1, $2)
RETURNING id, name, price, created_at;
