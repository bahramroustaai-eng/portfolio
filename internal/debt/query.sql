-- name: CreateDebt :one
INSERT INTO debts (lender, borrower, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDebtsByBorrowerID :many
SELECT * FROM debts
WHERE borrower = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountDebtsByBorrowerID :one
SELECT count(*) FROM debts
WHERE borrower = $1;