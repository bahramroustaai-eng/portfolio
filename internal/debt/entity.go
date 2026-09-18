package debt

import (
	"portfolio/internal/apperr"
	"time"
)

var (
	ErrPositiveAmount                = apperr.New(apperr.CodeInvalid, "amount must be greater than zero")
	ErrConflictDebtLenderAndBorrower = apperr.New(apperr.CodeConflict, "debt lender and borrower cannot be the same")
)

type DebtStatus string

const (
	Pending  DebtStatus = "pending"
	Paid     DebtStatus = "paid"
	Canceled DebtStatus = "canceled"
)

type Debt struct {
	ID               int32      `json:"id"`
	BorrowerID       int32      `json:"borrower_id"`
	BorrowerUsername string     `json:"borrower_username"`
	LenderID         int32      `json:"lender_id"`
	LenderUsername   string     `json:"lender_username"`
	Amount           int32      `json:"amount"`
	Status           DebtStatus `json:"status"`
}

type DebtPayment struct {
	ID         int32     `json:"id"`
	Amount     int64     `json:"amount"`
	DebtID     int64     `json:"debt_id"`
	PayerID    int64     `json:"payer_id"`
	ReceiverID int64     `json:"receiver_id"`
	Note       *string   `json:"note"`
	PaidAt     time.Time `json:"paid_at"`
	CreatedAt  time.Time `json:"created_at"`
}
