# Debt and Payment Feature Plan

## Goal

Allow users to record debts between users, view what they owe and what others owe them, record partial or full payments, and calculate the remaining balance safely.

The implementation should teach and demonstrate practical Go and PostgreSQL practices:

- clear package boundaries;
- request validation at the HTTP boundary;
- business rules in the service layer;
- database transactions for state changes;
- explicit SQL and deterministic results;
- authenticated ownership checks;
- immutable financial history.

## Scope

### Included

- Create a debt between two existing users.
- List debts where the current user is the borrower.
- List debts where the current user is the lender.
- Record a partial or full payment.
- Calculate paid and remaining amounts.
- Show open, partially paid, and paid states.
- Display pair balances on the user dashboard.

### Deferred

- Real bank or payment-provider integration.
- Payment refunds and reversals through an external provider.
- Splitting one payment across multiple debts.
- Notifications.
- Currency conversion.

## Domain rules

1. A debt has exactly one lender and one borrower.
2. The lender and borrower must be different users.
3. Debt and payment amounts must be positive.
4. A payment payer must be the debt borrower.
5. A payment receiver must be the debt lender.
6. A payment cannot exceed the current remaining debt.
7. A debt cannot be deleted after a payment exists.
8. The original debt amount never changes.
9. Remaining balance is derived from the original amount minus confirmed payments.
10. The authenticated user identity comes from the access token, not the request body.

## Proposed data model

### `debts`

Stores the original obligation.

Suggested fields:

```text
id
lender_id
borrower_id
amount
currency
description
due_at
status
created_at
```

Use `BIGINT` for money in the smallest currency unit. For the first version, use `IRR` and document whether amounts are rials or tomans. Do not use floating point values for money.

Suggested statuses:

```text
open
cancelled
```

`partially_paid` and `paid` can be derived from payment totals instead of stored as mutable state.

### `debt_payments`

Stores every payment event.

Suggested fields:

```text
id
debt_id
payer_id
receiver_id
amount
status
note
paid_at
created_at
```

For the first version, payment status can be:

```text
confirmed
```

Keep the status column so a future provider integration can add `pending`, `failed`, and `reversed` without redesigning the table.

## Balance calculation

For one debt:

```text
paid_amount = SUM(confirmed payments.amount)
remaining_amount = debt.amount - paid_amount
```

The API should return:

```json
{
  "id": 12,
  "lender_username": "ali",
  "borrower_username": "sara",
  "original_amount": 1000000,
  "paid_amount": 400000,
  "remaining_amount": 600000,
  "status": "partially_paid"
}
```

For a user pair:

```text
net_balance = amount_owed_to_current_user - amount_current_user_owes
```

Positive means the other user owes the current user. Negative means the current user owes the other user.

## Transaction design for recording a payment

The payment use case must run inside one database transaction:

1. Begin a transaction.
2. Lock the target debt row with `SELECT ... FOR UPDATE`.
3. Verify that the authenticated user is the borrower.
4. Calculate the confirmed paid amount.
5. Reject the request if the new payment exceeds the remaining amount.
6. Insert the payment.
7. Commit the transaction.

This protects against two concurrent payment requests overpaying the same debt.

The repository should accept a transaction-aware database handle or a transaction object. Do not call `context.Background()` inside repository methods when the caller already provides a request context.

## API proposal

```text
POST /api/v1/debts
GET  /api/v1/debts
GET  /api/v1/debts/{debt_id}
POST /api/v1/debts/{debt_id}/payments
GET  /api/v1/debts/{debt_id}/payments
GET  /api/v1/balances
```

### Create debt

The authenticated user is the lender. The request supplies only the borrower and debt details:

```json
{
  "borrower_id": 7,
  "amount": 1000000,
  "currency": "IRR",
  "description": "لپ‌تاپ"
}
```

### Create payment

The authenticated user must be the borrower. The request should not accept `payer_id` or `receiver_id`:

```json
{
  "amount": 400000,
  "note": "پرداخت بخشی از بدهی"
}
```

## Package responsibilities

### `internal/debt`

- domain types and status calculation;
- validation errors;
- debt and payment service methods;
- repository-facing models where useful.

### `internal/transport/http`

- decode request bodies;
- validate basic request shape;
- read authenticated user identity;
- map service errors to HTTP responses;
- never implement balance calculations directly.

### `internal/debt/db`

- generated sqlc code only;
- do not hand-edit generated files after the SQL workflow is established.

### `migrations`

- one migration for debt changes;
- one migration for payments;
- constraints and indexes belong in the database.

## Required indexes

Plan indexes for the main access patterns:

```text
debts(borrower_id, created_at DESC)
debts(lender_id, created_at DESC)
debt_payments(debt_id, status, paid_at)
```

## Acceptance criteria

- A user can create a debt only for another existing user.
- An unauthenticated request cannot create, list, or pay debts.
- A borrower sees debts they owe.
- A lender sees debts owed to them.
- A borrower can make a partial payment.
- A borrower can make the final payment and the debt becomes paid.
- Overpayment is rejected.
- Two concurrent payments cannot overpay a debt.
- Payment history is preserved.
- Empty list endpoints return `[]`, not `null`.
- All list queries have deterministic ordering.
- API responses never expose password fields.
- `go test ./...` passes.

## Learning review points

After implementation, review these areas:

1. Where is each business rule enforced?
2. Which rules are enforced by PostgreSQL constraints?
3. Which operations must share one transaction?
4. How does the code prevent a user from paying another user's debt?
5. What happens if the client retries a payment request?
6. Are repository methods using the request context?
7. Are generated sqlc files reproducible from the SQL source?
8. Can every balance be recomputed from stored records?
