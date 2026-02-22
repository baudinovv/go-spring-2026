package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"practice3/internal/_repository"
	"practice3/internal/handler"
	"practice3/internal/repository/_postgres"
	"practice3/internal/usecase"
	"practice3/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	pg := _postgres.NewPGXDialect(ctx, dbConfig)
	repos := repository.NewRepositories(pg)
	uc := usecase.NewUserUsecase(repos.UserRepository)
	router := handler.NewRouter(uc)

	// Middleware chain: logging -> auth -> router
	chain := handler.LoggingMiddleware(handler.AuthMiddleware(router))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      chain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:         "localhost",
		Port:         "5432",
		Username:     "postgres",
		Password:     "admin",
		DBName:       "mydb",
		SSLMode:      "disable",
		ExecTimeout:  5 * time.Second,
	}
}
