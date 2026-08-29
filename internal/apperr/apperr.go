package apperr

type Code int

const (
	CodeInvalid Code = iota
	CodeNotFound
	CodeConflict
	CodeUnauthorized
)

type Error struct {
	Code    Code
	Message string
}

func (e *Error) Error() string { return e.Message }

func New(code Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}
