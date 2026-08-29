-- name: CreateDebt :one
INSERT INTO debts (lender, borrower, amount)
VALUES ($1, $2, $3)
RETURNING *;