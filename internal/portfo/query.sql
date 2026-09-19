-- name: CreateItemType :one
INSERT INTO item_types(name)
VALUES ($1)
RETURNING *;

-- name: ListItemTypes :many
SELECT id, name
FROM item_types
ORDER BY name;

-- name: CreateItem :one
INSERT INTO items(user_id, name, type_id, total_cost, unit, price_per_unit, ticker, affect_profit, risk_level)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListItemsByUserID :many
SELECT id, user_id, name, type_id, total_cost, unit, ticker, affect_profit, risk_level, created_at, price_per_unit
FROM items
WHERE user_id = $1
ORDER BY created_at DESC;
