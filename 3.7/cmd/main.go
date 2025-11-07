package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"warehouse-control/internal/config"
	"warehouse-control/internal/handler"
	auth "warehouse-control/internal/infrastructure/jwt"
	"warehouse-control/internal/infrastructure/postgresql"
	"warehouse-control/internal/router"
	"warehouse-control/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := buildServer()
	var wg sync.WaitGroup

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
	DBConfig := config.GetDBConfig()
	db := postgresql.InitDB(DBConfig)

	itemRepo := postgresql.NewItemRepo(db)
	userRepo := postgresql.NewUserRepo(db)
	historyRepo := postgresql.NewHistoryRepo(db)

	jwt := auth.NewJWTProvider([]byte(config.GetJWTSecret()))

	authSvc := usecase.NewAuthService(jwt, userRepo)
	authHandler := handler.NewAuthHandler(authSvc, jwt)

	historySvc := usecase.NewHistoryService(historyRepo)
	historyHandler := handler.NewHistoryHandler(historySvc)

	svc := usecase.NewService(itemRepo)
	handler := handler.NewHandler(svc)

	router := router.NewRouter(handler, authHandler, historyHandler, jwt, db)

	return http.Server{
		Addr:    config.GetServerConfig(),
		Handler: router,
	}
}
