package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"eventbooker/internal/config"
	"eventbooker/internal/handler"
	"eventbooker/internal/infrastructure/postgres"
	"eventbooker/internal/router"
	"eventbooker/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var wg sync.WaitGroup

	srv := buildServer()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Listen and running :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server will shutdown gracefully...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down service: %v", err)
	}

	wg.Wait()
}

func buildServer() http.Server {
	dbConfig := config.LoadDBConfig()
	db, err := postgres.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("Error init database: %v", err)
	}

	svc := usecase.NewEventBookService(db)
	handler := handler.NewHandler(svc)
	router := router.NewEventBookRouter(handler)

	return http.Server{
		Addr:    config.LoadServerConfig(),
		Handler: router,
	}
}
