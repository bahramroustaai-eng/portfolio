-- name: CreateUser :one
INSERT INTO users (user_name, password)
VALUES ($1, $2)
RETURNING *;


-- name: GetUserByUsername :one
SELECT * FROM users
WHERE user_name = $1;


-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;