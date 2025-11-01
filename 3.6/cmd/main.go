package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"sales-tracker/internal/config"
	"sales-tracker/internal/handler"
	"sales-tracker/internal/infrastructure/postgresql"
	"sales-tracker/internal/router"
	"sales-tracker/internal/usecase"
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
	DBConfig := config.LoadDBConfig()
	db, err := postgresql.InitDB(DBConfig)
	if err != nil {
		log.Fatalf("error init DB: %v", err)
	}

	svc := usecase.NewTrackService(db)
	handler := handler.NewTrackHandler(svc)
	router := router.NewRouter(handler)

	return http.Server{
		Addr:    config.LoadSrvConfig(),
		Handler: router,
	}
}
