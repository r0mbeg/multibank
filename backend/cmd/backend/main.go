// cmd/backend/main.go

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
	"log/slog"
	_ "multibank/backend/docs"
	"multibank/backend/internal/app"
	"multibank/backend/internal/config"
	"multibank/backend/internal/logger"
	"multibank/backend/internal/shutdown"
	"os"
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
	shutdown.Wait(application, log)
}
