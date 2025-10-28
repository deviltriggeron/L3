package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"eventbooker/internal/domain"
)

type storage struct {
	db *sql.DB
}

func InitDB(cfg domain.DBConfig) (domain.Repo, error) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DB,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Database initialized successfully!")
	return &storage{db: db}, nil
}

func (s *storage) SelectEvent(id uuid.UUID) (*domain.EventBook, error) {
	var event domain.EventBook
	ctx := context.Background()

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
		SELECT e.event_id, e.date, e.event_info, u.user_id, u.user_name, e.seats_count, e.for_free, e.price, e.create_date
		FROM events e
		JOIN users u ON e.organizer_id = u.user_id
		WHERE e.event_id = $1;

	`
	err = tx.QueryRowContext(ctx, q, id).Scan(
		&event.EventID, &event.Date, &event.EventInfo,
		&event.Organizer.UserID, &event.Organizer.Name,
		&event.SeatsCount, &event.ForFree, &event.Price, &event.CreateDate,
	)

	if err != nil {
		return nil, err
	}
	err = tx.Commit()

	return &event, err
}

func (s *storage) SelectUser(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	ctx := context.Background()

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
		SELECT *
		FROM users
		WHERE user_id = $1
	`

	err = tx.QueryRowContext(ctx, q, id).Scan(
		&user.UserID, &user.Name,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	return &user, err
}

func (s *storage) InsertEvent(e *domain.EventBook) error {
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
        INSERT INTO events (
            event_id, date, event_info, organizer_id, seats_count, for_free, price
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.ExecContext(ctx, q, e.EventID, e.Date, e.EventInfo, e.Organizer.UserID, e.SeatsCount, e.ForFree, e.Price)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (s *storage) InsertUser(id uuid.UUID, name string) error {
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
		INSERT INTO users (
			user_id, user_name
		)
		VALUES ($1, $2);
	`

	_, err = tx.ExecContext(ctx, q, id, name)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (s *storage) InsertBooking(b *domain.Booking) error {
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
		INSERT INTO bookings(booking_id, event_id, user_id, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.ExecContext(ctx, q, b.BookingID, b.EventID, b.UserID, b.Status, b.CreatedAt, b.ExpiresAt)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (s *storage) UpdateBooking(eventID uuid.UUID, bookingID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
		UPDATE bookings SET
			status = $1
		WHERE event_id = $2 AND booking_id = $3
	`

	_, err = tx.ExecContext(ctx, q, domain.Paid, eventID, bookingID)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (s *storage) SelectBookingAt(eventID uuid.UUID, bookingID uuid.UUID) (bool, error) {
	var createdAt, expiredAt time.Time
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return false, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `
		SELECT created_at, expires_at
		FROM bookings
		WHERE event_id = $1 AND booking_id = $2
	`

	err = tx.QueryRowContext(ctx, q, eventID, bookingID).Scan(&createdAt, &expiredAt)
	if err != nil {
		log.Printf("error in query row %v", err)
		return false, err
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	fmt.Printf("Now : %v", time.Now())
	fmt.Printf("Expired time : %v", expiredAt)

	return time.Now().Before(expiredAt), nil
}

func (s *storage) CountPaidParticipants(eventID uuid.UUID) (int, error) {
	var count int
	q := `SELECT COUNT(*) FROM bookings WHERE event_id = $1 AND status = 'paid'`
	err := s.db.QueryRow(q, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *storage) PaidEvent(eventID uuid.UUID) (bool, error) {
	forFree := false
	var seatsCount int

	q := `
		SELECT for_free, seats_count
		FROM events
		WHERE event_id = $1
	`

	err := s.db.QueryRow(q, eventID).Scan(&forFree, &seatsCount)
	if err != nil {
		return false, err
	}

	participants, err := s.CountPaidParticipants(eventID)
	if err != nil {
		return false, fmt.Errorf("error count participants: %v", err)
	}

	if participants >= seatsCount {
		return forFree, fmt.Errorf("unfortunately, there are no more places for the event")
	}

	return forFree, nil
}
