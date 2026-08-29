package debt

import "portfolio/internal/apperr"

var (
	ErrPositiveAmount                = apperr.New(apperr.CodeInvalid, "amount must be greater than zero")
	ErrConflictDebtLenderAndBorrower = apperr.New(apperr.CodeConflict, "debt lender and borrower cannot be the same")
)

type Debt struct {
	ID       int32
	Borrower int32
	Amount   int32
	Lender   int32
}
