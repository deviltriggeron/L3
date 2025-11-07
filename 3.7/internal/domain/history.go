package domain

import (
	"time"

	"github.com/google/uuid"
)

type ItemHistory struct {
	ID        uuid.UUID
	ItemID    uuid.UUID
	Action    string
	OldData   map[string]any
	NewData   map[string]any
	ChangedBy string
	ChangedAt time.Time
}
