# Backlog

## Debt listing concurrency

**Status:** Backlog

Update `internal/debt.Service.ListDebts` to load the two independent debt groups concurrently:

- debts the authenticated user owes;
- debts owed to the authenticated user.

Use `errgroup.WithContext` so an error or request cancellation stops the other operation. Confirm that the PostgreSQL pool has enough available connections before enabling this. Measure the endpoint before and after the change so concurrency is justified by reduced latency.

Keep payment processing sequential inside its transaction. The row lock and the payment validation, insert, and status update must remain in one controlled transaction.

### Acceptance criteria

- Both debt groups are loaded concurrently with the request context.
- The first error is returned and the derived context is canceled.
- No data race is introduced.
- `go test -race ./...` passes.
- The endpoint latency is measured for sequential and concurrent implementations.
