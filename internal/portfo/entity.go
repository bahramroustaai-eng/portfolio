package portfo

import (
	"portfolio/internal/apperr"
	"strings"
	"time"
)

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

var (
	ErrInvalidItemTypeName = apperr.New(apperr.CodeInvalid, "item type name must not be empty")
	ErrItemTypeConflict    = apperr.New(apperr.CodeConflict, "item type already exists")
)

type ItemType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func NewItemTypeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrInvalidItemTypeName
	}
	return name, nil
}

type Item struct {
	ID             int32     `json:"id"`
	UserID         int32     `json:"user_id"`
	Name           string    `json:"name"`
	TypeID         int32     `json:"type_id"`
	TotalCost      int64     `json:"total_cost"`
	Unit           int32     `json:"unit"`
	PricePerUnit   int64     `json:"price_per_unit"`
	Ticker         string    `json:"ticker"`
	AffectedProfit bool      `json:"affected_profit"`
	RiskLevel      RiskLevel `json:"risk_level"`
	CreatedAt      time.Time `json:"created_at"`
}

func (i Item) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return apperr.New(apperr.CodeInvalid, "item name must not be empty")
	}
	if i.TypeID <= 0 {
		return apperr.New(apperr.CodeInvalid, "item type is required")
	}
	if i.TotalCost <= 0 {
		return apperr.New(apperr.CodeInvalid, "total cost must be greater than zero")
	}
	if i.Unit <= 0 {
		return apperr.New(apperr.CodeInvalid, "unit must be greater than zero")
	}
	if i.PricePerUnit <= 0 {
		return apperr.New(apperr.CodeInvalid, "price per unit must be greater than zero")
	}
	if i.RiskLevel != RiskLevelLow && i.RiskLevel != RiskLevelMedium && i.RiskLevel != RiskLevelHigh {
		return apperr.New(apperr.CodeInvalid, "invalid risk level")
	}
	return nil
}
