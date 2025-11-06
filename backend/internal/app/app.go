// internal/app/app.go

package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"multibank/backend/internal/logger"
	"multibank/backend/internal/service/bank"
	"multibank/backend/internal/service/product"
	"net/http"
	"strconv"
	"time"

	"multibank/backend/internal/config"
	httpserver "multibank/backend/internal/http-server"
	"multibank/backend/internal/service/auth"
	"multibank/backend/internal/service/user"
	"multibank/backend/internal/storage/sqlite"

	"multibank/backend/internal/service/auth/jwt"

	"multibank/backend/internal/service/openbanking"
)

type App struct {
	log     *slog.Logger
	cfg     *config.Config
	httpSrv *http.Server
	storage *sqlite.Storage
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	// --- storage (SQLite) ---
	st, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("sqlite init: %w", err)
	}

	// Migrations (db schema init)
	if err := st.Migrate(context.Background()); err != nil {
		_ = st.Close()
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}

	// --- repo + services ---
	userRepo := sqlite.NewUserRepo(st.DB())
	userSvc := user.New(log, userRepo)

	bankRepo := sqlite.NewBankRepo(st.DB())
	bankSvc := bank.New(log, bankRepo)

	productsClient := &openbanking.ProductsClient{
		HTTP: &http.Client{Timeout: 10 * time.Second},
	}
	prodSvc := product.New(log, bankRepo, bankSvc, productsClient)

	jwtMgr := jwt.New(cfg.HTTPServer.JWTSecret, cfg.HTTPServer.TokenTTL)
	authSvc := auth.New(log, userSvc, jwtMgr)

	// --- chi mux via httpserver.New ---
	srv := httpserver.New(
		httpserver.Deps{
			Logger:         log,
			UserService:    userSvc, // implements handlers.User
			AuthService:    authSvc, // implements handlers.Auth
			BankService:    bankSvc, // implements handlers.Bank
			ProductService: prodSvc, // implements handlers.Product
			JWT:            jwtMgr,
		},
		httpserver.Options{
			RequestTimeout:     cfg.HTTPServer.Timeout,
			BankEnsureOnStart:  true,             // check tokens on start
			BankEnsureInterval: 10 * time.Minute, // check tokens every ... minutes (e.g.)
		},
	)

	httpSrv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTPServer.Port),
		Handler:      srv.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: cfg.HTTPServer.Timeout + 2*time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		log:     log,
		cfg:     cfg,
		httpSrv: httpSrv,
		storage: st,
	}, nil
}

// Run — blocking HTTP-server start
func (a *App) Run() error {
	a.log.Info("http server starting",
		slog.String("addr", a.httpSrv.Addr),
		slog.String("env", a.cfg.Env),
		slog.String("storage_path", a.cfg.StoragePath),
		slog.String("log_level", a.cfg.Logger.LevelString),
	)

	// ListenAndServe blocks execution until the server shuts down
	if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("http server stopped with error", logger.Err(err))
	}
}

// Stop — graceful shutdown + resources closing
func (a *App) Stop(ctx context.Context) error {
	a.log.Info("shutting down http server...")
	if err := a.httpSrv.Shutdown(ctx); err != nil {
		a.log.Error("http shutdown error", logger.Err(err))
		return err
	}
	a.log.Info("closing storage...")
	if err := a.storage.Close(); err != nil {
		a.log.Error("storage close error", logger.Err(err))
		return err
	}
	a.log.Info("app stopped gracefully")
	return nil
}

// Export interfaces for server.New, if needed in tests.:
// var _ handlers.Auth = (*auth.Service)(nil)
// var _ handlers.User = (*user.Service)(nil)
//_ = handlers.Auth(nil)
//_ = handlers.User(nil)
