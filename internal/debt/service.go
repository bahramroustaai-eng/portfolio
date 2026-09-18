package debt

import (
	"context"
	"portfolio/internal/user"
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
}

func NewService(repo *Repository, userRepo *user.Repository) *Service {
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

	created, err := s.repo.CreateDebt(l.ID, b.ID, amount)
	if err != nil {
		return Debt{}, err
	}

	return created, nil
}

func (s *Service) ListDebt(ctx context.Context, userID int32, limit int32, offset int32) ([]Debt, int32, error) {
	debts, err := s.repo.GetDebtsByBorrowerID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.repo.CountDebtsByBorrowerID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return debts, totalCount, nil
}

func (s *Service) PayDebt(ctx context.Context, userID int32, debtID int32, amount int32, note *string) (DebtPayment, error) {
	if amount <= 0 {
		return DebtPayment{}, ErrPositiveAmount
	}
	return s.repo.PayDebt(ctx, userID, debtID, amount, note)
}
