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

type Repository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

var _ UserRepository = (*Repository)(nil)

func NewUserRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, q: db.New(pool)}
}

func (r *Repository) CreateUser(ctx context.Context, username string, password string) (User, error) {
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

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("could not get user: %w", err)
	}
	return toDomain(row), nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := r.q.GetUsers(ctx)
	if err != nil {
		return []User{}, fmt.Errorf("could not get users: %w", err)
	}
	users := make([]User, 0, len(rows))

	for _, row := range rows {
		users = append(users, toDomainGetUsers(row))
	}
	return users, nil
}

func toDomainGetUsers(i db.User) User {
	return User{
		ID:       i.ID,
		UserName: i.UserName,
	}
}

func toDomain(i db.User) User {
	return User{
		ID:       i.ID,
		UserName: i.UserName,
		Password: i.Password,
	}
}
