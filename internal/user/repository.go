package user

import (
	"context"
	"errors"
	"fmt"

	"portfolio/internal/user/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool, q: db.New(pool)}
}

func (r *UserRepository) CreateUser(ctx context.Context, username string, password string) (User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		UserName: username,
		Password: password,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrUserNameConflict
		}
		return User{}, fmt.Errorf("could not create user: %w", err)
	}
	return toDomain(row), nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("could not get user: %w", err)
	}
	return toDomain(row), nil
}

func toDomain(i db.User) User {
	return User{
		ID:       i.ID,
		UserName: i.UserName,
		Password: i.Password}
}
