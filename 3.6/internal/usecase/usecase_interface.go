package usecase

import (
	"github.com/google/uuid"

	"sales-tracker/internal/domain"
)

type TrackerService interface {
	Insert(item domain.Item) (uuid.UUID, error)
	Get(itemID string) (*domain.Item, error)
	Update(itemID string, item domain.Item) error
	Delete(itemID string) error
}

type TrackerServiceExtended interface {
	TrackerService
	GetAll(from string, to string, category string, typ string, limit string, offset string) ([]domain.Item, error)
	GetAnalytics(from string, to string, category string, typ string) (*domain.AnalyticsResult, error)
}
