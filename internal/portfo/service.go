package portfo

import "context"

type Service struct {
	repo ItemTypeRepository
}

func NewService(repo ItemTypeRepository) *Service {
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
