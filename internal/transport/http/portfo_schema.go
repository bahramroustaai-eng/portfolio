package transport

import (
	"portfolio/internal/portfo"
	"time"
)

type CreateItemTypeRequest struct {
	Name string `json:"name"`
}

type ItemTypeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CreateItemRequest struct {
	Name           string           `json:"name"`
	TypeID         int32            `json:"type_id"`
	TotalCost      int64            `json:"total_cost"`
	Unit           int32            `json:"unit"`
	PricePerUnit   int64            `json:"price_per_unit"`
	Ticker         string           `json:"ticker"`
	AffectedProfit bool             `json:"affected_profit"`
	RiskLevel      portfo.RiskLevel `json:"risk_level"`
}

type ItemResponse struct {
	ID             int32            `json:"id"`
	Name           string           `json:"name"`
	TypeID         int32            `json:"type_id"`
	TotalCost      int64            `json:"total_cost"`
	Unit           int32            `json:"unit"`
	PricePerUnit   int64            `json:"price_per_unit"`
	Ticker         string           `json:"ticker"`
	AffectedProfit bool             `json:"affected_profit"`
	RiskLevel      portfo.RiskLevel `json:"risk_level"`
	CreatedAt      time.Time        `json:"created_at"`
}

func toItemResponse(item portfo.Item) ItemResponse {
	return ItemResponse{
		ID:             item.ID,
		Name:           item.Name,
		TypeID:         item.TypeID,
		TotalCost:      item.TotalCost,
		Unit:           item.Unit,
		PricePerUnit:   item.PricePerUnit,
		Ticker:         item.Ticker,
		AffectedProfit: item.AffectedProfit,
		RiskLevel:      item.RiskLevel,
		CreatedAt:      item.CreatedAt,
	}
}
