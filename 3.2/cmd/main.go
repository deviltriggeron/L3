package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"shortener/internal/config"
	"shortener/internal/handler"
	"shortener/internal/router"
	"shortener/internal/service"
	db "shortener/internal/storage"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var wg sync.WaitGroup

	cfg := config.GetConfig()
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewShortenerService(dbConn)
	handler := handler.NewShortenerHandler(svc)
	router := router.NewRouter(handler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Listen and running :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	<-ctx.Done()
	srv.Shutdown(ctx)

	wg.Wait()
}
