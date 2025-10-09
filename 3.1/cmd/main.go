package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notifier/internal/handler"
	"notifier/internal/router"
	"notifier/internal/service"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var wg sync.WaitGroup

	svc := service.NewNotifierService()
	handler := handler.NewNotifierHandler(svc)
	router := router.NewRouter(*handler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Listen and running :8080\n")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)

	wg.Wait()
}
