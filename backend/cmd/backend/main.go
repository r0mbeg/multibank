package main

import (
	"log/slog"
	"multibank/backend/internal/config"
	"multibank/backend/internal/logger"
	"multibank/backend/internal/storage/sqlite"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// setup logger
	log := logger.Setup(cfg.Logger.Level)

	// connect to db
	st, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		return
	}
	defer st.Close()

	log.Info("starting multibank backend",
		slog.String("env", cfg.Env),
		slog.String("storage_path", cfg.StoragePath),
		slog.String("logger_level", cfg.Logger.LevelString),
		slog.Int("http_server_port", cfg.HTTPServer.Port),
		slog.String("http_server_timeout", cfg.HTTPServer.Timeout.String()),
		slog.String("http_server_token_ttl", cfg.HTTPServer.TokenTTL.String()),
	)

}
