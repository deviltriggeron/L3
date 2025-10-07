package entity

import "time"

const (
	Pending   = "pending"
	Sent      = "sent"
	Failed    = "failed"
	Cancelled = "cancelled"
)

type Notification struct {
	ID       int       `json:"id"`
	Message  string    `json:"message"`
	SendAt   time.Time `json:"send_at"`
	Status   string    `json:"status"`
	Attempts int       `json:"attempts"`
}
