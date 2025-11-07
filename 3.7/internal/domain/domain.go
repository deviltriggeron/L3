package domain

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID
	Product     string
	Price       float64
	Description string
	Count       int
	CreateDate  time.Time
}
