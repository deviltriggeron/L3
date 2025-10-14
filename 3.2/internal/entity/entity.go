package entity

import (
	"time"

	"github.com/google/uuid"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string
	ServerPort       string
}

type Transition struct {
	ID              uuid.UUID
	URLID           uuid.UUID
	UserAgent       string
	IPAddress       string
	TimeTransitions time.Time
}
