package entity

import "time"

const (
	Pending   = "pending"
	Sent      = "sent"
	Failed    = "failed"
	Cancelled = "cancelled"
)

type Notification struct {
	ID       int
	Message  string
	SendAt   time.Time
	Status   string
	Attempts int
}

type NotifierHandle struct {
	Message string    `json:"message"`
	SendAt  time.Time `json:"send_at"`
}
