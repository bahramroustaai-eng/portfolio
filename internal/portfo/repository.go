package portfo

import (
	"context"
	"errors"
	"fmt"
	"portfolio/internal/portfo/db"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemTypeRepository interface {
	CreateItemType(ctx context.Context, name string) (ItemType, error)
	ListItemTypes(ctx context.Context) ([]ItemType, error)
}

type PortfoRepository struct {
	q *db.Queries
}

var _ ItemTypeRepository = (*PortfoRepository)(nil)

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
