package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"sales-tracker/internal/domain"
)

func LoadDBConfig() domain.DBConfig {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("error load config DB: %v", err)
	}

	return domain.DBConfig{
		User:       os.Getenv("POSTGRES_USER"),
		Pass:       os.Getenv("POSTGRES_PASSWORD"),
		DB:         os.Getenv("POSTGRES_DB"),
		Host:       os.Getenv("POSTGRES_HOST"),
		Port:       os.Getenv("POSTGRES_PORT"),
		ServerPort: os.Getenv("POSTGRES_SERVER_PORT"),
	}
}

func LoadSrvConfig() string {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("error load server port: %v", err)
	}

	return os.Getenv("SERVER_PORT")
}
