package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"imageprocessor/internal/config"
	"imageprocessor/internal/domain"
	"imageprocessor/internal/handler"
	"imageprocessor/internal/infrastructure/kafka"
	"imageprocessor/internal/infrastructure/minio"
	"imageprocessor/internal/interfaces"
	"imageprocessor/internal/router"
	"imageprocessor/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var wg sync.WaitGroup

	minioRepo, broker, srv := buildServer()

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("Listen and running :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		startWorker(ctx, minioRepo, broker)
	}()

	<-ctx.Done()
	log.Println("Server will shutdown gracefully...")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutting down service: %v", err)
	}

	wg.Wait()
}

func buildServer() (interfaces.MinioRepository, interfaces.EventPublisher, *http.Server) {
	cfg := config.LoadConfigMinio()
	minioClient := minio.NewMinioRepo(*cfg)

	for i := 0; i < 10; i++ {
		err := minioClient.InitMinio()
		if err == nil {
			break
		}
		log.Printf("Waiting for Minio to be ready (%d/10)...", i+1)
		time.Sleep(5 * time.Second)
	}

	cfgBroker := config.LoadConfigKafka()
	broker := kafka.NewKafkaBroker(*cfgBroker)

	srvPort := config.LoadConfigServer()
	svc := usecase.NewImageProcService(minioClient, broker)
	handler := handler.NewImageProcHandler(svc)
	router := router.NewRouter(handler)

	return minioClient, broker, &http.Server{
		Addr:    srvPort,
		Handler: router,
	}
}

func startWorker(ctx context.Context, minioRepo interfaces.MinioRepository, kafkaBroker interfaces.EventPublisher) {
	kafkaBroker.Consume(ctx, func(key string, value []byte) error {
		var task domain.ImageTask
		if err := json.Unmarshal(value, &task); err != nil {
			return err
		}

		log.Printf("Processing task: %s, options: %v\n", task.ID, task.Options)
		return minioRepo.ProcessImage(ctx, task.ID, task.Options...)
	})
}
