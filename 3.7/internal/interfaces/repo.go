package interfaces

import (
	"github.com/google/uuid"

	"warehouse-control/internal/domain"
)

type ItemRepository interface {
	Select(itemID uuid.UUID) (*domain.Item, error)
	SelectAll() ([]domain.Item, error)
	Insert(item domain.Item) error
	Update(itemID uuid.UUID, item domain.Item) error
	Delete(itemID uuid.UUID) error
}

type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
}

type HistoryRepository interface {
	GetHistory() ([]domain.ItemHistory, error)
}
