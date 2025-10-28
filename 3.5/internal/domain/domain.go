package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID
	Name   string
}

type EventBook struct {
	EventID      uuid.UUID `json:"event_id"`
	Date         time.Time `json:"date"`
	EventInfo    string    `json:"event_info"`
	Organizer    User      `json:"organizer"`
	SeatsCount   int       `json:"seats_count"`
	Participants []User    `json:"participants"`
	ForFree      bool      `json:"for_free"`
	Price        float64   `json:"price"`
	CreateDate   time.Time `json:"create_date"`
}

type BookingStatus string

const (
	Pending  BookingStatus = "pending"
	Paid     BookingStatus = "paid"
	Canceled BookingStatus = "canceled"
)

type Booking struct {
	BookingID uuid.UUID     `json:"booking_id"`
	EventID   uuid.UUID     `json:"event_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Status    BookingStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	ExpiresAt time.Time     `json:"expires_at"`
}
