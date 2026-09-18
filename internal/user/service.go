package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewUserService(repo *Repository) *Service {
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

func (s *Service) Login(ctx context.Context, username, password string) (User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Service) GetUsersList(ctx context.Context) ([]User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return []User{}, err
	}
	return users, nil
}
