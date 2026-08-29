package transport

import (
	"encoding/json"
	"net/http"
	"portfolio/internal/debt"
	"strconv"
)

type DebtHandler struct {
	svc *debt.Service
}

func NewDebtHandler(svc *debt.Service) *DebtHandler {
	return &DebtHandler{svc: svc}
}

func (h *DebtHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/debt", h.CreateDebt)
	mux.Handle("GET /api/v1/debts", RequireAuth(http.HandlerFunc(h.GetDebts)))
}

type CreateDebtRequest struct {
	Lender   string `json:"lender"`
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
	var req CreateDebtRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	created, err := h.svc.CreateDebt(r.Context(), req.Lender, req.Borrower, req.Amount)
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
		totalAmount += d.Amount
		debtByLender[d.LenderUsername] += d.Amount
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
