# portfolio

Go service for tracking users, debt, and investment portfolios, with notification delivery.

## Stack

- Go 1.26
- `net/http` stdlib router (`ServeMux`)
- PostgreSQL via `pgx/v5`
- `sqlc` for generated query code
- `goose` for migrations
- `docker-compose` for local Postgres

## Project layout

Vertical-slice-per-domain: each business domain owns its entity, repository, SQL queries, and generated DB code. This mirrors the existing `notification` package — new domains should follow the same shape rather than introducing a separate layered (`domain/usecase/repo`) architecture.

```
portfolio/
├── cmd/
│   └── api/
│       └── main.go                    # composition root: config, db pool, repos, handlers, router, server
├── internal/
│   ├── user/
│   │   ├── entity.go                  # domain type
│   │   ├── repository.go              # interface + Postgres impl
│   │   ├── query.sql                  # sqlc source
│   │   └── db/                        # sqlc generated (db.go, models.go, query.sql.go)
│   ├── debt/
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── query.sql
│   │   └── db/
│   ├── portfolio/
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── query.sql
│   │   └── db/
│   ├── notification/                  # existing, keep as reference pattern
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── query.sql
│   │   └── db/
│   ├── platform/
│   │   ├── config/                    # env parsing, out of main.go
│   │   └── postgres/                  # pool builder (move newPool here)
│   └── transport/
│       └── http/
│           ├── router.go              # mounts all domain routes
│           ├── middleware.go
│           ├── user_handler.go        # one handler file per domain
│           ├── debt_handler.go
│           ├── portfolio_handler.go
│           └── notification_handler.go
├── migrations/                        # goose migrations, timestamp-prefixed
├── sqlc.yaml                          # one entry per domain
├── docker-compose.yml
├── Makefile
├── go.mod / go.sum
```

## Domain mapping

| Business concept | Workflow | API | DB entity | Code |
|---|---|---|---|---|
| User | signup / login | `POST /api/v1/users`, `POST /api/v1/login` | `users` | `internal/user/*` |
| Debt | create / track / pay debt | `POST /api/v1/debts`, `GET /api/v1/debts` | `debts` | `internal/debt/*` |
| Portfolio | asset allocation, valuation | `GET /api/v1/portfolio` | `portfolios`, `holdings` | `internal/portfolio/*` |
| Notification | alert delivery | `POST /api/v1/notification` | `notifications` | `internal/notification/*` |

## Known gaps / cleanup

- [ ] `sqlc.yaml` only has a `sql:` entry for `notification` — add one per domain (`user`, `debt`, `portfolio`).
- [ ] `user` package has `query.sql` but no `entity.go`, `repository.go`, or generated `db/`.
- [ ] `debt` and `portfolio` packages, migrations, and tables don't exist yet.
- [ ] `main.go` builds the notification repo but discards it (`_ = notification.NewPostgresRepository(...)`) instead of injecting into the handler.
- [ ] `Handler` in `transport/http` is a single stub struct with no dependencies — split into one handler per domain, each constructed with its own repo.
- [ ] Primary key inconsistency: `users` migration uses `serial`, notifications use `uuid` (via `google/uuid` + sqlc override). Standardize on `uuid` for new tables.

## Local development

```
make up             # start Postgres (docker-compose)
make migrate-up      # apply migrations (goose)
make sqlc            # regenerate sqlc code
make test            # go test ./... -race
```
