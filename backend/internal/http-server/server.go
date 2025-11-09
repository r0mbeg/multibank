// internal/http-server/server.go

package httpserver

import (
	"context"
	"multibank/backend/internal/logger"
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
	mux      *chi.Mux
	logger   *slog.Logger
	shutdown chan struct{} // channel for graceful completing background cycles
}

type Deps struct {
	Logger             *slog.Logger
	UserService        handlers.User
	AuthService        handlers.Auth
	BankService        handlers.Bank
	ProductService     handlers.Product
	RecommendedService handlers.Recommended
	ConsentService     handlers.Consent
	AccountService     handlers.Account
	JWT                *jwt.Manager
}

type Options struct {
	RequestTimeout time.Duration

	BankEnsureOnStart  bool          // getting tokens on startup
	BankEnsureInterval time.Duration // periodic ensure, 0 = disable
	BankEnsureWorkers  int

	ConsentEnsureOnStart  bool
	ConsentEnsureInterval time.Duration
	ConsentEnsureWorkers  int
}

func New(deps Deps, opts Options) *Server {
	r := chi.NewRouter()

	// CORS settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // cache preflight in seconds
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

	// Protected routes /products
	r.Route("/products", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterProductRoutes(rr, deps.ProductService)
	})

	// Protected routes /recommended-products
	r.Route("/admin/recommended-products", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		// TODO: protect for Admin
		handlers.RegisterRecommendedRoutes(rr, deps.RecommendedService)
	})

	// Protected routes /consents
	r.Route("/consents", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterConsentRoutes(rr, deps.ConsentService)
	})

	// Protected routes /accounts
	r.Route("/accounts", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterAccountRoutes(rr, deps.AccountService) // передай в Deps
	})

	// swagger ui
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	srv := &Server{
		mux:      r,
		logger:   deps.Logger,
		shutdown: make(chan struct{}),
	}

	// ensure Banks TOKENS
	if opts.BankEnsureOnStart {
		go srv.ensureBanksOnce(deps, 20*time.Second, opts.BankEnsureWorkers)
	}
	if opts.BankEnsureInterval > 0 {
		go srv.runBankEnsureLoop(deps, opts)
	}

	// ensure Account Consents
	if opts.ConsentEnsureOnStart || opts.ConsentEnsureInterval > 0 {
		go srv.runConsentEnsureLoop(deps, opts)
	}
	return srv
}

func (s *Server) Handler() stdhttp.Handler { return s.mux }

// Stop stops background loops (should be used in graceful shutdown)
func (s *Server) Stop() {
	select {
	case <-s.shutdown:
		// already closed
	default:
		close(s.shutdown)
	}
}

// ensureBanksOnce
func (s *Server) ensureBanksOnce(deps Deps, timeout time.Duration, workers int) {
	if workers <= 0 {
		workers = 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := deps.BankService.EnsureTokensForEnabledWithWorkers(ctx, workers); err != nil {
		s.logger.Warn("ensure bank tokens on start failed",
			slog.Any("err", err),
			slog.Int("workers", workers),
		)
	} else {
		s.logger.Info("ensure bank tokens on start ok",
			slog.Int("workers", workers),
		)
	}
}

// runBankEnsureLoop periodically checks bank tokens
func (s *Server) runBankEnsureLoop(deps Deps, opt Options) {
	if opt.BankEnsureInterval <= 0 {
		return
	}
	t := time.NewTicker(opt.BankEnsureInterval)
	defer t.Stop()

	for {
		select {
		case <-s.shutdown:
			s.logger.Info("stopping bank ensure loop")
			return
		case <-t.C:
			workers := opt.BankEnsureWorkers
			if workers <= 0 {
				workers = 1
			}
			// даём половину интервала на выполнение цикла
			ctx, cancel := context.WithTimeout(context.Background(), opt.BankEnsureInterval/2)
			if err := deps.BankService.EnsureTokensForEnabledWithWorkers(ctx, workers); err != nil {
				s.logger.Warn("scheduled bank token ensure failed",
					slog.Any("err", err),
					slog.Int("workers", workers),
				)
			} else {
				s.logger.Info("scheduled bank token ensure ok",
					slog.Int("workers", workers),
				)
			}
			cancel()
		}
	}
}

func (s *Server) runConsentEnsureLoop(deps Deps, opt Options) {
	log := s.logger.With(slog.String("component", "consent-ensure"))
	workers := opt.ConsentEnsureWorkers
	if workers <= 0 {
		workers = 2
	}

	// разово на старте
	if opt.ConsentEnsureOnStart {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		n, err := deps.ConsentService.RefreshStale(ctx, 100, workers)
		cancel()
		if err != nil {
			log.Warn("initial consent ensure failed", logger.Err(err))
		} else {
			log.Info("initial consent ensure done", slog.Int("updated", n))
		}
	}

	if opt.ConsentEnsureInterval <= 0 {
		return
	}

	ticker := time.NewTicker(opt.ConsentEnsureInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdown:
			log.Info("stopping consent ensure loop")
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), opt.ConsentEnsureInterval/2)
			n, err := deps.ConsentService.RefreshStale(ctx, 100, workers)
			cancel()
			if err != nil {
				log.Warn("periodic consent ensure failed", logger.Err(err))
			} else if n > 0 {
				log.Info("periodic consent ensure", slog.Int("updated", n))
			}
		}
	}
}
