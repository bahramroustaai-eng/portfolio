package transport

import (
	"portfolio/internal/debt"
	"time"
)

type CreateDebtRequest struct {
	Borrower string `json:"borrower"`
	Amount   int32  `json:"amount"`
}

type CreateDebtResponse struct {
	ID int32 `json:"id"`
}

type DebtGroupResponse struct {
	Debts       []debt.Debt      `json:"debts"`
	DebtByUser  map[string]int32 `json:"debt_by_user"`
	TotalAmount int32            `json:"total_amount"`
	TotalCount  int32            `json:"total_count"`
}

type GetDebtsResponse struct {
	IOwe     DebtGroupResponse `json:"i_owe"`
	OwedToMe DebtGroupResponse `json:"owed_to_me"`
	Limit    int32             `json:"limit"`
	Offset   int32             `json:"offset"`
}

type PayDebtRequest struct {
	Amount int32   `json:"amount"`
	Note   *string `json:"note"`
}

type DebtPaymentResponse struct {
	ID         int32     `json:"id"`
	Amount     int32     `json:"amount"`
	DebtID     int32     `json:"debt_id"`
	PayerID    int32     `json:"payer_id"`
	ReceiverID int32     `json:"receiver_id"`
	Note       *string   `json:"note"`
	PaidAt     time.Time `json:"paid_at"`
	CreatedAt  time.Time `json:"created_at"`
}
