package debt

import (
	"context"
	"portfolio/internal/user"
)

type Service struct {
	repo     *Repository
	userRepo *user.UserRepository
}

func NewService(repo *Repository, userRepo *user.UserRepository) *Service {
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
