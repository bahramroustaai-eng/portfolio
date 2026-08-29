package user

import (
	"time"

	"portfolio/internal/apperr"
)

var (
	ErrUserNameConflict   = apperr.New(apperr.CodeConflict, "user name already exists")
	ErrNotFound           = apperr.New(apperr.CodeNotFound, "user not found")
	ErrInvalidInput       = apperr.New(apperr.CodeInvalid, "invalid input")
	ErrInvalidCredentials = apperr.New(apperr.CodeUnauthorized, "invalid username or password")
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
