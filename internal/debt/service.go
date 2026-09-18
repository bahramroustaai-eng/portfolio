package debt

import (
	"context"
	"portfolio/internal/user"
)

type DebtRepository interface {
	CreateDebt(ctx context.Context, lenderID, borrowerID, amount int32) (Debt, error)
	GetDebtsByBorrowerID(ctx context.Context, borrowerID, limit, offset int32) ([]Debt, error)
	GetDebtsByLenderID(ctx context.Context, lenderID, limit, offset int32) ([]Debt, error)
	CountDebtsByBorrowerID(ctx context.Context, borrowerID int32) (int32, error)
	CountDebtsByLenderID(ctx context.Context, lenderID int32) (int32, error)
	PayDebt(ctx context.Context, payerID, debtID, amount int32, note *string) (DebtPayment, error)
}

type UserReader interface {
	GetUserByUsername(ctx context.Context, username string) (user.User, error)
}

type DebtList struct {
	Debts      []Debt
	TotalCount int32
}

type DebtLists struct {
	IOwe     DebtList
	OwedToMe DebtList
}

type Service struct {
	repo     DebtRepository
	userRepo UserReader
}

func NewDebtService(repo DebtRepository, userRepo UserReader) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateDebt(ctx context.Context, lender string, borrower string, amount int32) (Debt, error) {

	if amount <= 0 {
		return Debt{}, ErrPositiveAmount
	}

	l, err := s.userRepo.GetUserByUsername(ctx, lender)
	if err != nil {
		return Debt{}, err
	}
	b, err := s.userRepo.GetUserByUsername(ctx, borrower)
	if err != nil {
		return Debt{}, err
	}
	if l.ID == b.ID {
		return Debt{}, ErrConflictDebtLenderAndBorrower
	}

	created, err := s.repo.CreateDebt(ctx, l.ID, b.ID, amount)
	if err != nil {
		return Debt{}, err
	}

	return created, nil
}

func (s *Service) ListDebts(ctx context.Context, userID int32, limit int32, offset int32) (DebtLists, error) {
	iOwe, err := s.listDebts(
		ctx,
		userID,
		limit,
		offset,
		s.repo.GetDebtsByBorrowerID,
		s.repo.CountDebtsByBorrowerID,
	)
	if err != nil {
		return DebtLists{}, err
	}

	owedToMe, err := s.listDebts(
		ctx,
		userID,
		limit,
		offset,
		s.repo.GetDebtsByLenderID,
		s.repo.CountDebtsByLenderID,
	)
	if err != nil {
		return DebtLists{}, err
	}

	return DebtLists{IOwe: iOwe, OwedToMe: owedToMe}, nil
}

type listDebtsFunc func(context.Context, int32, int32, int32) ([]Debt, error)
type countDebtsFunc func(context.Context, int32) (int32, error)

func (s *Service) listDebts(ctx context.Context, userID, limit, offset int32, list listDebtsFunc, count countDebtsFunc) (DebtList, error) {
	debts, err := list(ctx, userID, limit, offset)
	if err != nil {
		return DebtList{}, err
	}
	totalCount, err := count(ctx, userID)
	if err != nil {
		return DebtList{}, err
	}
	return DebtList{Debts: debts, TotalCount: totalCount}, nil
}

func (s *Service) PayDebt(ctx context.Context, userID int32, debtID int32, amount int32, note *string) (DebtPayment, error) {
	if amount <= 0 {
		return DebtPayment{}, ErrPositiveAmount
	}
	return s.repo.PayDebt(ctx, userID, debtID, amount, note)
}
