package httpserver

import (
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

	// базовые middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mwLogger.New(deps.Logger)) // твой middleware, пишущий через slog
	r.Use(middleware.Recoverer)

	if opts.RequestTimeout > 0 {
		r.Use(middleware.Timeout(opts.RequestTimeout))
	}

	// Публичные маршруты (регистрация/логин)
	handlers.RegisterAuthRoutes(r, handlers.AuthDeps{
		Auth: deps.AuthService,
		JWT:  deps.JWT,
	})

	// Защищённые маршруты /users/*
	r.Route("/users", func(rr chi.Router) {
		rr.Use(authmw.Auth(deps.JWT))
		handlers.RegisterUserRoutes(rr, deps.UserService)
	})

	return &Server{mux: r}
}

func (s *Server) Handler() stdhttp.Handler { return s.mux }
