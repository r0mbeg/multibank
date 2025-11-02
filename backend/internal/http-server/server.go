package httpserver

import (
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/cors"
	"log/slog"
	"multibank/backend/internal/auth"
	"multibank/backend/internal/auth/jwt"
	authmw "multibank/backend/internal/auth/middleware"
	"multibank/backend/internal/http-server/handlers"
	mwLogger "multibank/backend/internal/http-server/middleware/logger"
	"multibank/backend/internal/service/user"
)

type Server struct {
	mux *chi.Mux
}

type Deps struct {
	Logger      *slog.Logger
	UserService *user.Service
	AuthService *auth.Service
	JWT         *jwt.Manager
}

type Options struct {
	RequestTimeout time.Duration
}

func New(deps Deps, opts Options) *Server {
	r := chi.NewRouter()

	// CORS settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // если отправляешь куки/авторизацию через fetch(..., credentials:'include')
		MaxAge:           300,  // кэш preflight в секундах
	}))

	// базовые middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mwLogger.New(deps.Logger)) // middleware with metadata of requests
	r.Use(middleware.Recoverer)

	if opts.RequestTimeout > 0 {
		r.Use(middleware.Timeout(opts.RequestTimeout))
	}

	// Public routes (registration/login)
	handlers.RegisterAuthRoutes(r, handlers.AuthDeps{
		Auth: deps.AuthService,
		JWT:  deps.JWT,
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

	// swagger ui
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return &Server{mux: r}
}

func (s *Server) Handler() stdhttp.Handler { return s.mux }
