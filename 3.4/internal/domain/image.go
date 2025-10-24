package domain

import "time"

type ImageStatus string

const (
	StatusPending   ImageStatus = "pending"
	StatusProcessed ImageStatus = "processed"
	StatusFailed    ImageStatus = "failed"
)

type Image struct {
	ID          string      `json:"id"`
	FileName    string      `json:"file_name"`
	ContentType string      `json:"content_type"`
	Size        int64       `json:"size"`
	Status      ImageStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ImageProcessOption int

const (
	ProcessResize ImageProcessOption = iota
	ProcessThumbnail
	ProcessWatermark
)
