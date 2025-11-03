// @title           Multibank API
// @version         1.0
// @description     API для аутентификации и пользователей.
// @BasePath        /
// @schemes         http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log/slog"
	_ "multibank/backend/docs"
	"multibank/backend/internal/app"
	"multibank/backend/internal/config"
	"multibank/backend/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// go run ./cmd/backend --config=./config/local.yaml

func main() {

	// load config
	cfg := config.MustLoad()

	// setup logger
	log := logger.Setup(cfg.Logger.Level)

	// app building
	application, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to init app", slog.Any("err", err))
		os.Exit(1)
	}

	// app start
	go application.MustRun()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	log.Info("received shutdown signal", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Stop(ctx); err != nil {
		log.Error("graceful shutdown failed", slog.Any("err", err))
	}
}
