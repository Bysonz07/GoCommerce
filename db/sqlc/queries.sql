-- name: AddToCart :exec
INSERT INTO carts (user_id, product_id, quantity) VALUES ($1, $2, $3);

-- name: GetCartByUser :many
SELECT * FROM carts WHERE user_id = $1;

-- name: CreateOrder :exec
INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3);

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = $1;