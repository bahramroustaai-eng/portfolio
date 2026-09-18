-- name: CreateDebt :one
INSERT INTO debts (lender_id, borrower_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDebtsByBorrowerID :many
SELECT d.id,
       d.amount,
       COALESCE(SUM(p.amount), 0)::int AS paid_amount,
       d.amount - COALESCE(SUM(p.amount), 0)::int AS remaining_amount,
       d.created_at,
       d.status,
       d.lender_id,
       lu.user_name AS lender_username,
       d.borrower_id,
       bu.user_name AS borrower_username
FROM debts d
         LEFT JOIN debt_payments p ON p.debt_id = d.id
         JOIN users lu ON lu.id = d.lender_id
         JOIN users bu ON bu.id = d.borrower_id
WHERE d.borrower_id = $1
GROUP BY d.id, d.amount, d.created_at, d.status, d.lender_id, lu.user_name, d.borrower_id, bu.user_name
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountDebtsByBorrowerID :one
SELECT count(*)
FROM debts
WHERE borrower_id = $1;

-- name: GetDebtByID :one
SELECT *
FROM debts
WHERE id = $1;

-- name: GetDebtForUpdate :one
SELECT *
FROM debts
WHERE id = $1
FOR UPDATE;

-- name: GetPaidAmountByDebtID :one
SELECT COALESCE(SUM(amount), 0)::int
FROM debt_payments
WHERE debt_id = $1;

-- name: CreateDebtPayment :one
INSERT INTO debt_payments (payer_id, receiver_id, debt_id, amount, note)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: MarkDebtPaid :exec
UPDATE debts
SET status = 'paid'
WHERE id = $1;
