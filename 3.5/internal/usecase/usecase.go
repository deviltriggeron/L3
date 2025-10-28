package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"eventbooker/internal/domain"
)

type eventBookService struct {
	storage domain.Repo
}

func NewEventBookService(db domain.Repo) EventBookService {
	return &eventBookService{
		storage: db,
	}
}

func (s *eventBookService) CreateEvent(e domain.EventBook) error {
	if s.checkUserInDB(e.Organizer) {
		err := s.storage.InsertEvent(&e)
		if err != nil {
			return fmt.Errorf("error create event: %v", err)
		}
	} else {
		return fmt.Errorf("user not registred in app")
	}
	return nil
}

func (s *eventBookService) CreateUser(user domain.User) error {
	if err := s.storage.InsertUser(user.UserID, user.Name); err != nil {
		return fmt.Errorf("error create user: %v", err)
	}

	return nil
}

func (s *eventBookService) Booking(user string, event string) (uuid.UUID, error) {
	userID, err := uuid.Parse(user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parse user id: %v", err)
	}
	eventID, err := uuid.Parse(event)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parse event id: %v", err)
	}

	status, err := s.checkEvent(eventID)
	if err != nil {
		return uuid.Nil, err
	}

	if !status {
		bookingID := uuid.New()

		booking := &domain.Booking{
			BookingID: bookingID,
			EventID:   eventID,
			UserID:    userID,
			Status:    domain.Pending,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}

		err = s.storage.InsertBooking(booking)
		if err != nil {
			return uuid.Nil, fmt.Errorf("error create booking: %v", err)
		}

		return bookingID, nil
	} else {
		bookingID := uuid.New()

		booking := &domain.Booking{
			BookingID: bookingID,
			EventID:   eventID,
			UserID:    userID,
			Status:    domain.Paid,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
		}

		err = s.storage.InsertBooking(booking)
		if err != nil {
			return uuid.Nil, fmt.Errorf("error create booking for free event: %v", err)
		}

		return uuid.Nil, nil
	}
}

func (s *eventBookService) Payment(event string, booking string) error {
	eventID, err := uuid.Parse(event)
	if err != nil {
		return fmt.Errorf("error parse event id: %v", err)
	}
	bookingID, err := uuid.Parse(booking)
	if err != nil {
		return fmt.Errorf("error parse booking id: %v", err)
	}

	ok := s.checkBookingInDB(eventID, bookingID)
	if !ok {
		return fmt.Errorf("payment time has passed")
	}

	err = s.storage.UpdateBooking(eventID, bookingID)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventBookService) GetEvent(stringID string) (*domain.EventBook, error) {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return nil, fmt.Errorf("error parse id: %v", err)
	}

	event, err := s.storage.SelectEvent(id)
	if err != nil {
		return nil, fmt.Errorf("get event from DB: %v", err)
	}

	count, err := s.storage.CountPaidParticipants(id)
	if err != nil {
		return nil, fmt.Errorf("error count participants: %v", err)
	}

	event.Participants = make([]domain.User, count)

	return event, nil
}

func (s *eventBookService) checkUserInDB(user domain.User) bool {
	_, err := s.storage.SelectUser(user.UserID)
	return err == nil
}

func (s *eventBookService) checkBookingInDB(eventID uuid.UUID, bookingID uuid.UUID) bool {
	status, err := s.storage.SelectBookingAt(eventID, bookingID)
	if err != nil {
		log.Printf("error check booking: %v", err)
		return false
	}

	return status
}

func (s *eventBookService) checkEvent(eventID uuid.UUID) (bool, error) {
	return s.storage.PaidEvent(eventID)
}
