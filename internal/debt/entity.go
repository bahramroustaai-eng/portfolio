package debt

import "portfolio/internal/apperr"

var (
	ErrPositiveAmount                = apperr.New(apperr.CodeInvalid, "amount must be greater than zero")
	ErrConflictDebtLenderAndBorrower = apperr.New(apperr.CodeConflict, "debt lender and borrower cannot be the same")
)

type Debt struct {
	ID               int32  `json:"id"`
	BorrowerID       int32  `json:"borrower_id"`
	BorrowerUsername string `json:"borrower_username"`
	LenderID         int32  `json:"lender_id"`
	LenderUsername   string `json:"lender_username"`
	Amount           int32  `json:"amount"`
}
