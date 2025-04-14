package main

import (
	"assignment/internal/config"
	"assignment/internal/libs/logger"
	httpserver "assignment/internal/server"

	"assignment/internal/repository/postgres"
	"assignment/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envPath        = ".env"
	migrationsPath = "./migrations"
)

func main() {
	ctx := context.Background()
	var (
		logger     = must(logger.New())
		cfg        = must(config.New(envPath))
		repository = must(postgres.New(ctx, cfg.PGConfig.DSN, cfg.PassSecret))
	)

	defer logger.Sync()
	defer repository.DB.Close()

	if os.Getenv("ENV") != "prod" {
		migrations, err := migrate.New(
			fmt.Sprintf("file://%s", migrationsPath),
			cfg.PGConfig.DSN,
		)
		if err != nil {
			logger.Fatal(err)
		}

		if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
			logger.Fatal(err)
		}
	}

	svc := service.New(repository, cfg.JWTSecret, logger)
	srv := httpserver.New(cfg, logger, svc)

	go func() {
		logger.Info("server started on port: ", cfg.HTTPServer.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server shutdown: %s", err)
	}

	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}

	logger.Info("server exiting")
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}
