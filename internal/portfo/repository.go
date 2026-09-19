package portfo

import (
	"context"
	"errors"
	"fmt"
	"portfolio/internal/apperr"
	"portfolio/internal/portfo/db"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PortfoRepository struct {
	q *db.Queries
}

var _ Repository = (*PortfoRepository)(nil)

func NewRepository(pool *pgxpool.Pool) *PortfoRepository {
	return &PortfoRepository{q: db.New(pool)}
}

func (r *PortfoRepository) CreateItemType(ctx context.Context, name string) (ItemType, error) {
	row, err := r.q.CreateItemType(ctx, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ItemType{}, ErrItemTypeConflict
		}
		return ItemType{}, fmt.Errorf("create item type: %w", err)
	}
	return ItemType{ID: row.ID, Name: row.Name}, nil
}

func (r *PortfoRepository) ListItemTypes(ctx context.Context) ([]ItemType, error) {
	rows, err := r.q.ListItemTypes(ctx)
	if err != nil {
		return nil, err
	}

	itemTypes := make([]ItemType, 0, len(rows))
	for _, row := range rows {
		itemTypes = append(itemTypes, ItemType{ID: row.ID, Name: row.Name})
	}
	return itemTypes, nil
}

func (r *PortfoRepository) CreateItem(
	ctx context.Context,
	userID int32,
	params Item,
) (Item, error) {
	row, err := r.q.CreateItem(ctx, db.CreateItemParams{
		UserID:       userID,
		Name:         params.Name,
		TypeID:       &params.TypeID,
		TotalCost:    params.TotalCost,
		Unit:         params.Unit,
		PricePerUnit: params.PricePerUnit,
		Ticker:       &params.Ticker,
		AffectProfit: &params.AffectedProfit,
		RiskLevel:    db.RiskLevel(params.RiskLevel),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return Item{}, apperr.New(apperr.CodeNotFound, "item type not found")
		}
		return Item{}, fmt.Errorf("create item: %w", err)
	}
	return toDomainItem(row.ID, row.UserID, row.Name, row.TypeID, row.TotalCost, row.Unit, row.Ticker, row.AffectProfit, row.RiskLevel, row.CreatedAt, row.PricePerUnit), nil
}

func (r *PortfoRepository) ListItemsByUserID(ctx context.Context, userID int32) ([]Item, error) {
	rows, err := r.q.ListItemsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	items := make([]Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainItem(row.ID, row.UserID, row.Name, row.TypeID, row.TotalCost, row.Unit, row.Ticker, row.AffectProfit, row.RiskLevel, row.CreatedAt, row.PricePerUnit))
	}
	return items, nil
}

func toDomainItem(id, userID int32, name string, typeID *int32, totalCost int64, unit int32, ticker *string, affectProfit *bool, riskLevel db.RiskLevel, createdAt time.Time, pricePerUnit int64) Item {
	item := Item{
		ID:           id,
		UserID:       userID,
		Name:         name,
		TotalCost:    totalCost,
		Unit:         unit,
		PricePerUnit: pricePerUnit,
		RiskLevel:    RiskLevel(riskLevel),
		CreatedAt:    createdAt,
	}
	if typeID != nil {
		item.TypeID = *typeID
	}
	if ticker != nil {
		item.Ticker = *ticker
	}
	if affectProfit != nil {
		item.AffectedProfit = *affectProfit
	}
	return item
}
