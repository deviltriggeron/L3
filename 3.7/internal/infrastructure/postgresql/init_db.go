package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"warehouse-control/internal/domain"
)

func InitDB(cfg domain.DBconfig) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DB)
	db, err := sql.Open(cfg.DB, dsn)
	if err != nil {
		log.Fatalf("error open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error ping DB: %v", err)
	}

	fmt.Println("Database initialized successfully!")
	return db
}
