// internal/http-server/server.go

package httpserver

import (
	"context"
	"multibank/backend/internal/service/auth/jwt"
	authmw "multibank/backend/internal/service/auth/middleware"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"log/slog"
	"multibank/backend/internal/http-server/handlers"
	mwLogger "multibank/backend/internal/http-server/middleware/logger"

	"github.com/go-chi/cors"
)

type Server struct {
	mux *chi.Mux
}

type Deps struct {
	Logger      *slog.Logger
	UserService handlers.User
	AuthService handlers.Auth
	BankService handlers.Bank
	JWT         *jwt.Manager
}

type Options struct {
	RequestTimeout     time.Duration
	BankEnsureOnStart  bool          // getting tokens on startup
	BankEnsureInterval time.Duration // periodic ensure, 0 = disable
}

func New(deps Deps, opts Options) *Server {
	r := chi.NewRouter()

	// CORS settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // if you send cookies/authorization via fetch(..., credentials:'include')
		MaxAge:           300,  // cache preflight in seconds
	}))

	// basic middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mwLogger.New(deps.Logger)) // middleware with metadata of requests
	r.Use(middleware.Recoverer)

	if opts.RequestTimeout > 0 {
		r.Use(middleware.Timeout(opts.RequestTimeout))
	}

	// Public routes (registration/login)
	/*
		handlers.RegisterAuthRoutes(r, handlers.AuthDeps{
			Auth: deps.AuthService,
			JWT:  deps.JWT,
		})*/
	r.Route("/auth", func(rr chi.Router) {
		// rr.Use(authmw.Auth(deps.JWT)) DO NOT NEED !!!
		handlers.RegisterAuthRoutes(rr, deps.AuthService, deps.JWT)
	})

	// Protected routes /users/*
	r.Route("/users", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterUserRoutes(rr, deps.UserService)
	})

	// Protected routes /me/*
	r.Route("/me", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterMeRoutes(rr, deps.UserService)
	})

	// Protected routes /banks
	r.Route("/banks", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterBankRoutes(rr, deps.BankService)
	})

	// swagger ui
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	if opts.BankEnsureOnStart {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			if err := deps.BankService.EnsureTokensForEnabled(ctx); err != nil {
				deps.Logger.Warn("ensure tokens on start failed", slog.Any("err", err))
			}
		}()
	}

	if opts.BankEnsureInterval > 0 {
		go func() {
			t := time.NewTicker(opts.BankEnsureInterval)
			defer t.Stop()
			for range t.C {
				// no timeout (~but can be small)
				if err := deps.BankService.EnsureTokensForEnabled(context.Background()); err != nil {
					deps.Logger.Warn("scheduled token ensure failed", slog.Any("err", err))
				}
			}
		}()
	}

	return &Server{mux: r}
}

func (s *Server) Handler() stdhttp.Handler { return s.mux }
