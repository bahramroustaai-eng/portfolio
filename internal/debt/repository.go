package debt

import (
	"context"
	"fmt"
	"portfolio/internal/debt/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, q: db.New(pool)}
}

func (r *Repository) CreateDebt(lender int32, borrower int32, amount int32) (Debt, error) {
	debt, err := r.q.CreateDebt(context.Background(), db.CreateDebtParams{
		Lender:   lender,
		Borrower: borrower,
		Amount:   amount,
	})
	if err != nil {
		return Debt{}, fmt.Errorf("create debt: %w", err)
	}
	return r.toDomain(debt), nil
}

func (r *Repository) toDomain(debt db.Debt) Debt {
	return Debt{
		ID:       debt.ID,
		Lender:   debt.Lender,
		Borrower: debt.Borrower,
		Amount:   debt.Amount,
	}
}
