package debt

import (
	"context"
	"errors"
	"fmt"
	"portfolio/internal/debt/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

var _ DebtRepository = (*Repository)(nil)

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, q: db.New(pool)}
}

func (r *Repository) CreateDebt(ctx context.Context, lender int32, borrower int32, amount int32) (Debt, error) {
	debt, err := r.q.CreateDebt(ctx, db.CreateDebtParams{
		LenderID:   lender,
		BorrowerID: borrower,
		Amount:     amount,
	})
	if err != nil {
		return Debt{}, fmt.Errorf("create debt: %w", err)
	}
	return r.toDomain(debt), nil
}

func (r *Repository) GetDebtsByBorrowerID(ctx context.Context, borrowerID int32, limit int32, offset int32) ([]Debt, error) {
	rows, err := r.q.GetDebtsByBorrowerID(ctx, db.GetDebtsByBorrowerIDParams{BorrowerID: borrowerID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("get debts: %w", err)
	}
	debts := make([]Debt, 0, len(rows))
	for _, row := range rows {
		debts = append(debts, Debt{
			ID:               row.ID,
			LenderID:         row.LenderID,
			LenderUsername:   row.LenderUsername,
			BorrowerUsername: row.BorrowerUsername,
			BorrowerID:       row.BorrowerID,
			Amount:           row.Amount,
			PaidAmount:       row.PaidAmount,
			RemainingAmount:  row.RemainingAmount,
			Status:           DebtStatus(row.Status),
		})
	}
	return debts, nil
}

func (r *Repository) CountDebtsByBorrowerID(ctx context.Context, borrowerID int32) (int32, error) {
	count, err := r.q.CountDebtsByBorrowerID(ctx, borrowerID)
	if err != nil {
		return 0, fmt.Errorf("count debts: %w", err)
	}
	return int32(count), nil
}

func (r *Repository) toDomain(debt db.Debt) Debt {
	return Debt{
		ID:         debt.ID,
		LenderID:   debt.LenderID,
		BorrowerID: debt.BorrowerID,
		Amount:     debt.Amount,
		Status:     DebtStatus(debt.Status),
	}
}

func (r *Repository) GetDebtByID(ctx context.Context, debtID int32) (Debt, error) {
	row, err := r.q.GetDebtByID(ctx, debtID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Debt{}, ErrDebtNotFound
		}
		return Debt{}, fmt.Errorf("get debt: %w", err)
	}
	return r.toDomain(row), nil
}

func (r *Repository) PayDebt(ctx context.Context, payerID, debtID, amount int32, note *string) (DebtPayment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return DebtPayment{}, fmt.Errorf("begin payment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	debt, err := q.GetDebtForUpdate(ctx, debtID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DebtPayment{}, ErrDebtNotFound
		}
		return DebtPayment{}, fmt.Errorf("get debt for payment: %w", err)
	}
	if debt.BorrowerID != payerID {
		return DebtPayment{}, ErrDebtPaymentUnauthorized
	}
	if DebtStatus(debt.Status) == Canceled {
		return DebtPayment{}, ErrDebtCanceled
	}

	paidAmount, err := q.GetPaidAmountByDebtID(ctx, debtID)
	if err != nil {
		return DebtPayment{}, fmt.Errorf("get paid amount: %w", err)
	}
	remainingAmount := debt.Amount - paidAmount
	if amount > remainingAmount {
		return DebtPayment{}, ErrPaymentExceedsRemaining
	}

	payment, err := q.CreateDebtPayment(ctx, db.CreateDebtPaymentParams{
		PayerID:    payerID,
		ReceiverID: debt.LenderID,
		DebtID:     debtID,
		Amount:     amount,
		Note:       note,
	})
	if err != nil {
		return DebtPayment{}, fmt.Errorf("create debt payment: %w", err)
	}
	if amount == remainingAmount {
		if err := q.MarkDebtPaid(ctx, debtID); err != nil {
			return DebtPayment{}, fmt.Errorf("mark debt paid: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return DebtPayment{}, fmt.Errorf("commit payment transaction: %w", err)
	}

	return DebtPayment{
		ID:         payment.ID,
		Amount:     payment.Amount,
		PayerID:    payment.PayerID,
		ReceiverID: payment.ReceiverID,
		Note:       payment.Note,
		DebtID:     payment.DebtID,
		PaidAt:     payment.PaidAt,
		CreatedAt:  payment.CreatedAt,
	}, nil
}
