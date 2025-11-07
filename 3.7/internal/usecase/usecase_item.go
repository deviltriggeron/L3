package usecase

import (
	"github.com/google/uuid"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

type itemService struct {
	storage interfaces.ItemRepository
}

func NewService(db interfaces.ItemRepository) ControllerService {
	return &itemService{
		storage: db,
	}
}

func (s *itemService) AddItem(item domain.Item) (uuid.UUID, error) {
	id := uuid.New()

	item.ID = id

	err := s.storage.Insert(item)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *itemService) GetItem(stringID string) (*domain.Item, error) {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return nil, err
	}

	return s.storage.Select(id)
}

func (s *itemService) GetAllItem() ([]domain.Item, error) {
	return s.storage.SelectAll()
}

func (s *itemService) UpdateItem(stringID string, item domain.Item) error {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return err
	}

	return s.storage.Update(id, item)
}

func (s *itemService) DeleteItem(stringID string) error {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return err
	}

	return s.storage.Delete(id)
}
