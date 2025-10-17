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

type Comments struct {
	CommentID uuid.UUID  `db:"comment_id"`
	ParentID  *uuid.UUID `db:"parent_id"`
	UserName  string     `db:"user_name"`
	Comment   string     `db:"comment"`
	Path      string     `db:"path"`
	Date      time.Time  `db:"date"`
}

type CommentResponse struct {
	UserName string
	Comment  string
	ParentID uuid.UUID
}
