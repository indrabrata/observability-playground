-- name: GetProduct :one
SELECT * FROM products
WHERE id = ? LIMIT 1;

-- name: GetProducts :many
SELECT * FROM products
ORDER BY name;

-- name: CreateProduct :one
INSERT INTO products (
  name, quantity, price, created_at
) VALUES (
  ?, ?, ?, ?
) RETURNING id, name, quantity, price, created_at, updated_at;

-- name: UpdateProduct :exec
UPDATE products
set name = ?,
quantity = ?,
price = ?,
updated_at = ?
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = ?;
