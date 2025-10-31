package main

import (
	"context"
	"log/slog"
	"multibank/backend/internal/config"
	httpserver "multibank/backend/internal/http-server"
	"multibank/backend/internal/logger"
	"multibank/backend/internal/service"
	"multibank/backend/internal/storage/sqlite"
	"net/http"
	"strconv"
	"time"
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

	ctx := context.Background()

	if err := st.Migrate(ctx); err != nil {
		log.Error("migrate", "err", err)
		return
	}

	// repo + service
	repo := sqlite.NewUserRepo(st.DB())
	svc := service.New(repo, cfg.HTTPServer.Timeout)

	srv := httpserver.New(svc, cfg.HTTPServer.Timeout)
	httpSrv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTPServer.Port),
		Handler:      srv.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: cfg.HTTPServer.Timeout + 2*time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("http server starting", "addr", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("http server error", "err", err)
	}
}
