package transport

import (
	"encoding/json"
	"net/http"
	"portfolio/internal/debt"
	"strconv"
	"strings"
	"time"
)

type DebtHandler struct {
	svc *debt.Service
}

func NewDebtHandler(svc *debt.Service) *DebtHandler {
	return &DebtHandler{svc: svc}
}

func (h *DebtHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/debts", RequireAuth(http.HandlerFunc(h.CreateDebt)))
	mux.Handle("GET /api/v1/debts", RequireAuth(http.HandlerFunc(h.GetDebts)))
	mux.Handle("POST /api/v1/debts/{debt_id}/payments", RequireAuth(http.HandlerFunc(h.PayDebt)))
}

type CreateDebtRequest struct {
	Borrower string `json:"borrower"`
	Amount   int32  `json:"amount"`
}

type CreateDebtResponse struct {
	ID int32 `json:"id"`
}

// CreateDebt godoc
//
//	@Summary		Create a debt
//	@Description	Records a debt between a lender and a borrower
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateDebtRequest	true	"debt payload"
//	@Success		201		{object}	CreateDebtResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/debt [post]
func (h *DebtHandler) CreateDebt(w http.ResponseWriter, r *http.Request) {
	userName, ok := UserNameFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	var req CreateDebtRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	created, err := h.svc.CreateDebt(r.Context(), userName, req.Borrower, req.Amount)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, CreateDebtResponse{ID: created.ID})

}

// GetDebts godoc
//
//	@Summary		List debts
//	@Description	Lists debts where the authenticated user is the borrower
//	@Tags			debts
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	GetDebtsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/debts [get]
func (h *DebtHandler) GetDebts(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	limit := int32(defaultLimit)
	if v := r.URL.Query().Get("limit"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil || parsed <= 0 {
			writeError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = int32(parsed)
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	offset := int32(0)
	if v := r.URL.Query().Get("offset"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil || parsed < 0 {
			writeError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = int32(parsed)
	}

	debts, totalCount, err := h.svc.ListDebt(r.Context(), userID, limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	debtByLender := make(map[string]int32)
	var totalAmount int32
	for _, d := range debts {
		totalAmount += d.RemainingAmount
		debtByLender[d.LenderUsername] += d.RemainingAmount
	}

	writeJSON(w, http.StatusOK, GetDebtsResponse{
		Debts:        debts,
		TotalAmount:  totalAmount,
		DebtByLender: debtByLender,
		Limit:        limit,
		Offset:       offset,
		TotalCount:   totalCount,
	})
}

type GetDebtsResponse struct {
	Debts        []debt.Debt      `json:"debts"`
	DebtByLender map[string]int32 `json:"debt_by_lender"`
	TotalAmount  int32            `json:"total_amount"`
	Limit        int32            `json:"limit"`
	Offset       int32            `json:"offset"`
	TotalCount   int32            `json:"total_count"`
}

// PayDebt godoc
//
//	@Summary		Record a debt payment
//	@Description	Creates a partial or full payment for the specified debt.
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			debt_id	path	int32	true	"Debt ID"
//	@Param			request	body	PayDebtRequest	true	"payment payload"
//	@Success		201	{object}	DebtPaymentResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/debts/{debt_id}/payments [post]
func (h *DebtHandler) PayDebt(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized user")
		return
	}

	debtID := strings.TrimSpace(r.PathValue("debt_id"))

	parsedDebtID, err := strconv.ParseInt(debtID, 10, 32)
	if err != nil || parsedDebtID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid debt id")
		return
	}

	var req PayDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	debtPayment, err := h.svc.PayDebt(r.Context(), userID, int32(parsedDebtID), req.Amount, req.Note)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, DebtPaymentResponse(debtPayment))
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
