package domain

type ImageTask struct {
	ID      string               `json:"id"`
	Options []ImageProcessOption `json:"options"`
}
