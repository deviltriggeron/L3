package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"imageprocessor/internal/domain"
)

func LoadConfigMinio() *domain.MinioCfg {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ssl := os.Getenv("MINIO_USE_SSL") == "1" || os.Getenv("MINIO_USE_SSL") == "true"

	return &domain.MinioCfg{
		Endpoint:        os.Getenv("MINIO_ENDPOINT"),
		Bucket:          os.Getenv("MINIO_BUCKET"),
		BucketProcessed: os.Getenv("MINIO_BUCKET_PROC"),
		AccessKey:       os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:       os.Getenv("MINIO_SECRET_KEY"),
		SSL:             ssl,
	}
}

func LoadConfigServer() string {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return os.Getenv("SERVER_PORT")
}

func LoadConfigKafka() *domain.ConfigBroker {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &domain.ConfigBroker{
		Broker:  os.Getenv("KAFKA_BROKER"),
		GroupID: os.Getenv("KAFKA_GROUP"),
		Topic:   os.Getenv("KAFKA_TOPIC"),
	}
}
