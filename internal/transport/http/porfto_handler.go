package transport

import (
	"net/http"
	"portfolio/internal/user"
)

type PortfoHandler struct {
	svc *user.Service
}

func NewPortfoHandler(svc *user.Service) *PortfoHandler {
	return &PortfoHandler{svc: svc}
}

func (h *PortfoHandler) RegisterRoutes(mux *http.ServeMux) {

}

func (h *PortfoHandler) Item(w http.ResponseWriter, r *http.Request) {}
