package user

import (
	"errors"
	"time"
)

var (
	ErrUserNameConflict   = errors.New("user name already exists")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type User struct {
	ID        int32
	UserName  string
	Password  string
	CreatedAt time.Time
}

type CreateInput struct {
	UserName string
	Password string
}
