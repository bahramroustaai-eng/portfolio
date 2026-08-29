-- name: CreateDebt :one
INSERT INTO debts (lender_id, borrower_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDebtsByBorrowerID :many
SELECT
    d.id, d.amount, d.created_at,
    d.lender_id, lu.user_name AS lender_username,
    d.borrower_id, bu.user_name AS borrower_username
FROM debts d
         JOIN users lu ON lu.id = d.lender_id
         JOIN users bu ON bu.id = d.borrower_id
WHERE d.borrower_id = $1
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountDebtsByBorrowerID :one
SELECT count(*) FROM debts
WHERE borrower_id = $1;