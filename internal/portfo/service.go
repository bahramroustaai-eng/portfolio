package portfo

import "context"

type Repository interface {
	CreateItemType(ctx context.Context, name string) (ItemType, error)
	ListItemTypes(ctx context.Context) ([]ItemType, error)
	CreateItem(ctx context.Context, userID int32, item Item) (Item, error)
	ListItemsByUserID(ctx context.Context, userID int32) ([]Item, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateItemType(ctx context.Context, name string) (ItemType, error) {
	cleanName, err := NewItemTypeName(name)
	if err != nil {
		return ItemType{}, err
	}

	created, err := s.repo.CreateItemType(ctx, cleanName)
	if err != nil {
		return ItemType{}, err
	}
	return created, nil
}

func (s *Service) ListItemTypes(ctx context.Context) ([]ItemType, error) {
	return s.repo.ListItemTypes(ctx)
}

func (s *Service) CreateItem(ctx context.Context, userID int32, item Item) (Item, error) {
	if err := item.Validate(); err != nil {
		return Item{}, err
	}

	created, err := s.repo.CreateItem(ctx, userID, item)
	if err != nil {
		return Item{}, err
	}
	return created, nil
}

func (s *Service) ListItemsByUserID(ctx context.Context, userID int32) ([]Item, error) {
	return s.repo.ListItemsByUserID(ctx, userID)
}
