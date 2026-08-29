package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, input CreateInput) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("could not hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, input.UserName, string(hash))
	if err != nil {
		return User{}, err
	}
	return user, nil
}
