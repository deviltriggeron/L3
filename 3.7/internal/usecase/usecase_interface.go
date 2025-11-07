package usecase

import (
	"github.com/google/uuid"

	"warehouse-control/internal/domain"
)

type ControllerService interface {
	AddItem(item domain.Item) (uuid.UUID, error)
	GetItem(stringID string) (*domain.Item, error)
	GetAllItem() ([]domain.Item, error)
	UpdateItem(stringID string, item domain.Item) error
	DeleteItem(stringID string) error
}

type AuthService interface {
	Login(username string, password string) (string, error)
}

type HistoryService interface {
	GetHistory() ([]domain.ItemHistory, error)
}
