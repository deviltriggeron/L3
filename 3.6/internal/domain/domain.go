package domain

import (
	"time"

	"github.com/google/uuid"
)

type ItemType string

const (
	Income  ItemType = "income"
	Expense ItemType = "expense"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	Type        ItemType  `json:"type"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type AnalyticsResult struct {
	Sum          float64 `json:"sum"`
	Avg          float64 `json:"avg"`
	Count        int     `json:"count"`
	Median       float64 `json:"median"`
	Percentile90 float64 `json:"percentile_90"`
}

type ItemFilter struct {
	From     string
	To       string
	Category string
	Type     ItemType
	Limit    int
	Offset   int
}

type AnalyticsFilter struct {
	From     string
	To       string
	Category string
	Type     ItemType
}
