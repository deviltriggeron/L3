package usecase

import (
	"eventbooker/internal/domain"

	"github.com/google/uuid"
)

type EventBookService interface {
	CreateEvent(e domain.EventBook) error
	CreateUser(user domain.User) error
	Booking(user string, event string) (uuid.UUID, error)
	Payment(user string, event string) error
	GetEvent(stringID string) (*domain.EventBook, error)
}
