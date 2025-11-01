package domain

import (
	"github.com/google/uuid"
)

type StorageRepository interface {
	Select(itemID uuid.UUID) (*Item, error)
	Insert(item Item) error
	Update(itemID uuid.UUID, item Item) error
	Delete(itemID uuid.UUID) error
}

type StorageRepositoryExtended interface {
	StorageRepository
	GetAll(filter ItemFilter) ([]Item, error)
	GetAnalytics(filter AnalyticsFilter) (*AnalyticsResult, error)
}
