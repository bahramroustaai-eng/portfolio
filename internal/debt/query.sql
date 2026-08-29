-- name: CreateDebt :one
INSERT INTO debts (lender, borrower, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDebtsByBorrowerID :many
SELECT * FROM debts
WHERE borrower = $1;