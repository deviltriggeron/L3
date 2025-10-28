package domain

import "github.com/google/uuid"

type Repo interface {
	SelectEvent(id uuid.UUID) (*EventBook, error)
	SelectUser(id uuid.UUID) (*User, error)
	InsertEvent(e *EventBook) error
	InsertUser(id uuid.UUID, name string) error
	InsertBooking(booking *Booking) error
	UpdateBooking(eventID uuid.UUID, bookingID uuid.UUID) error
	SelectBookingAt(eventID uuid.UUID, bookingID uuid.UUID) (bool, error)
	CountPaidParticipants(eventID uuid.UUID) (int, error)
	PaidEvent(eventID uuid.UUID) (bool, error)
}
