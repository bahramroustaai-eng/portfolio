package transport

import (
	"encoding/json"
	"net/http"
	"portfolio/internal/portfo"
)

type PortfolioHandler struct {
	svc *portfo.Service
}

func NewPortfolioHandler(svc *portfo.Service) *PortfolioHandler {
	return &PortfolioHandler{svc: svc}
}

func (h *PortfolioHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /item-types", RequireAuth(http.HandlerFunc(h.ListItemTypes)))
	mux.Handle("POST /item-types", RequireAuth(http.HandlerFunc(h.CreateItemType)))
}

// CreateItemType godoc
//
//	@Summary		Create an investment item type
//	@Tags			investment item types
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	CreateItemTypeRequest	true	"item type payload"
//	@Success		201	{object}	ItemTypeResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Router			/item-types [post]
func (h *PortfolioHandler) CreateItemType(w http.ResponseWriter, r *http.Request) {
	var request CreateItemTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()
	createItemType, err := h.svc.CreateItemType(r.Context(), request.Name)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, ItemTypeResponse{ID: createItemType.ID, Name: createItemType.Name})
}

// ListItemTypes godoc
//
//	@Summary		List investment item types
//	@Tags			investment item types
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	ItemTypeResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/item-types [get]
func (h *PortfolioHandler) ListItemTypes(w http.ResponseWriter, r *http.Request) {
	itemTypes, err := h.svc.ListItemTypes(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]ItemTypeResponse, 0, len(itemTypes))
	for _, itemType := range itemTypes {
		response = append(response, ItemTypeResponse{ID: itemType.ID, Name: itemType.Name})
	}
	writeJSON(w, http.StatusOK, response)
}
