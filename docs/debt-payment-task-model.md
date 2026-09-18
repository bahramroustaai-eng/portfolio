# Debt and Payment Task Model

This document breaks the feature into small implementation tasks. Complete the tasks in order and keep each change reviewable.

## Task states

Use these states for personal progress:

```text
BACKLOG -> READY -> IN_PROGRESS -> REVIEW -> DONE
                         |
                         v
                      BLOCKED
```

### State meaning

- `BACKLOG`: identified but not prepared.
- `READY`: requirements and dependencies are clear.
- `IN_PROGRESS`: actively implementing.
- `REVIEW`: implementation is complete and needs tests or code review.
- `BLOCKED`: progress needs a design decision or external dependency.
- `DONE`: acceptance criteria and verification are complete.

## Task template

Every task should record:

```text
Title:
Goal:
Dependencies:
Files or packages:
Database changes:
API changes:
Business rules:
Tests:
Manual verification:
Learning notes:
Definition of done:
```

## Work breakdown

### D001 — Confirm money and payment semantics

Goal: decide whether amounts are rials or tomans and whether the first payment flow is a manual confirmation flow.

Dependencies: none.

Definition of done:

- currency unit is documented;
- payment status behavior is documented;
- partial payment behavior is documented;
- deferred provider integration is explicitly out of scope.

Learning focus: domain modeling and avoiding ambiguous money representations.

### D002 — Design the debt migration

Goal: extend the debt model with the fields required by the agreed design.

Files:

- `migrations/`;
- `internal/debt/query.sql`;
- `sqlc.yaml` if the query configuration needs correction.

Definition of done:

- foreign keys are present;
- self-debts are rejected;
- positive amounts are enforced;
- currency representation is explicit;
- useful indexes are present;
- migration has a working down path where the project convention supports it.

Learning focus: PostgreSQL constraints, indexes, and migration safety.

### D003 — Add the payment migration

Goal: create the immutable payment history table.

Definition of done:

- payment references a debt;
- payer and receiver reference users;
- amount is positive;
- payment status is constrained;
- payment timestamps are present;
- debt lookup index exists.

Learning focus: relational integrity and event/history tables.

### D004 — Regenerate sqlc code

Goal: make generated database code reproducible from SQL files.

Definition of done:

- `sqlc generate` completes;
- generated files are formatted;
- generated output is not manually customized;
- queries return explicit columns instead of `SELECT *` where practical.

Learning focus: code generation and keeping source of truth clear.

### D005 — Implement debt detail and balance queries

Goal: return original, paid, remaining, and derived status values.

Definition of done:

- borrower and lender views use the correct user ID;
- paid totals include only confirmed payments;
- empty results return an empty slice;
- ordering is deterministic;
- passwords are never selected into API response models.

Learning focus: SQL aggregation, nullable values, and DTO boundaries.

### D006 — Implement transactional payment service

Goal: record a payment safely.

Definition of done:

- payment use case accepts request context;
- debt row is locked before remaining balance is checked;
- payer is verified against the authenticated borrower;
- overpayment returns a domain error;
- payment insert and validation commit atomically;
- rollback behavior is tested.

Learning focus: transactions, row locks, race conditions, and service boundaries.

### D007 — Add payment HTTP endpoints

Goal: expose payment creation and payment history through authenticated routes.

Definition of done:

- invalid JSON returns `400`;
- missing or invalid authentication returns `401`;
- unknown debt returns `404`;
- unauthorized borrower access returns an appropriate error;
- overpayment returns a stable domain error;
- success responses have documented shapes.

Learning focus: HTTP contracts and error mapping.

### D008 — Add user balance endpoint

Goal: expose net balances between the current user and other users.

Definition of done:

- balances include both debts created by the user and debts owed to the user;
- confirmed payments reduce the correct direction;
- zero balances are handled consistently;
- result ordering is deterministic.

Learning focus: query design and separating directional balances from net balances.

### D009 — Update the Persian RTL UI

Goal: show debts owed, debts receivable, remaining amounts, and a payment action.

Definition of done:

- current user's identity comes from the session;
- user selection cannot choose the current user;
- payment form sends only the amount and optional note;
- loading, empty, success, and error states are translated;
- paid debts cannot be paid again;
- UI never trusts a client-side balance as the final authority.

Learning focus: frontend/API contracts and defensive client behavior.

### D010 — Add tests and run the verification gate

Goal: verify the feature as a complete workflow.

Required checks:

- `go test ./...`;
- service test for self-debt rejection;
- service test for partial payment;
- service test for exact final payment;
- service test for overpayment;
- authorization test for another user's debt;
- repository or integration test for concurrent payment behavior if the test setup supports PostgreSQL;
- `git diff --check`.

Learning focus: testing business rules instead of only testing HTTP handlers.

## Review checklist before marking DONE

- Is the database the source of truth for balances?
- Are all payment state changes transactional?
- Can a request be safely retried?
- Are authorization checks based on authenticated identity?
- Are SQL queries explicit and ordered?
- Are errors represented as stable domain errors?
- Are generated files regenerated instead of manually edited?
- Does the implementation preserve payment history?
- Does the API return consistent empty arrays?
- Are tests proving the important invariants?

