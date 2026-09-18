-- name: CreateItemType :one
INSERT INTO item_types(name)
VALUES ($1)
RETURNING *;

-- name: ListItemTypes :many
SELECT id, name
FROM item_types
ORDER BY name;
