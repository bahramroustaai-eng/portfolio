package transport

import (
	"encoding/json"
	"net/http"
	"portfolio/internal/debt"
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
//	@Success		200	{array}		debt.Debt
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/debts [get]
func (h *DebtHandler) GetDebts(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	debts, err := h.svc.ListDebt(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, debts)
}
