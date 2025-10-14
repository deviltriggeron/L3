package config

import (
	"log"
	"os"
	e "shortener/internal/entity"

	"github.com/joho/godotenv"
)

func GetConfig() *e.Config {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &e.Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		ServerPort:       os.Getenv("SERVER_PORT"),
	}
}
