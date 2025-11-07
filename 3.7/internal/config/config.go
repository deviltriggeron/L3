package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"warehouse-control/internal/domain"
)

func GetDBConfig() domain.DBconfig {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("error load config")
	}

	return domain.DBconfig{
		User:       os.Getenv("POSTGRES_USER"),
		Pass:       os.Getenv("POSTGRES_PASSWORD"),
		DB:         os.Getenv("POSTGRES_DB"),
		Host:       os.Getenv("POSTGRES_HOST"),
		Port:       os.Getenv("POSTGRES_PORT"),
		ServerPort: os.Getenv("POSTGRES_SERVER_PORT"),
	}
}

func GetServerConfig() string {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("error load config")
	}

	return os.Getenv("SERVER_PORT")
}

func GetJWTSecret() string {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("error load config")
	}

	return os.Getenv("SECRET_JWT")
}
