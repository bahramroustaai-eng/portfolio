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

-- name: GetUsers :many
SELECT id, user_name, password, created_at
FROM users
ORDER BY user_name;
