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
	mux.Handle("GET /items", RequireAuth(http.HandlerFunc(h.ListItems)))
	mux.Handle("POST /items", RequireAuth(http.HandlerFunc(h.CreateItem)))
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

// CreateItem godoc
//
//	@Summary		Create new item
//	@Tags			investment items
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	CreateItemRequest	true	"investment item payload"
//	@Success		201	{object}	ItemResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Router			/items [post]
func (h *PortfolioHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	var request CreateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	newItem, err := h.svc.CreateItem(r.Context(), userID, portfo.Item{
		Name:           request.Name,
		TypeID:         request.TypeID,
		TotalCost:      request.TotalCost,
		Unit:           request.Unit,
		PricePerUnit:   request.PricePerUnit,
		Ticker:         request.Ticker,
		AffectedProfit: request.AffectedProfit,
		RiskLevel:      request.RiskLevel,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toItemResponse(newItem))
}

// ListItems godoc
//
//	@Summary		List investment items
//	@Tags			investment items
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	ItemResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items [get]
func (h *PortfolioHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	items, err := h.svc.ListItemsByUserID(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toItemResponse(item))
	}
	writeJSON(w, http.StatusOK, response)
	return
}
