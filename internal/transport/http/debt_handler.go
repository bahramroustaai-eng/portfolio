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
}

type CreateDebtRequest struct {
	Lender   string `json:"lender"`
	Borrower string `json:"borrower"`
	Amount   int32  `json:"amount"`
}

type CreateDebtResponse struct {
	ID int32 `json:"id"`
}

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
