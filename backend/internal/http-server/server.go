package httpserver

import (
	"log/slog"
	"multibank/backend/internal/auth"
	"multibank/backend/internal/auth/jwt"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	r.Use(mwLogger.New(deps.Logger))
	r.Use(middleware.Recoverer)

	// таймаут запросов — один раз на уровне HTTP
	if opts.RequestTimeout > 0 {
		r.Use(middleware.Timeout(opts.RequestTimeout))
	}

	// Регистрация feature-роутов (каждая фича сама вешает свои пути)
	//handlers.RegisterUserRoutes(r, deps.UserService)

	// Registration auth routes
	handlers.RegisterAuthRoutes(r, handlers.AuthDeps{
		Auth: deps.AuthService,
		JWT:  deps.JWT,
	})
	handlers.RegisterUserRoutes(r, deps.UserService)

	return &Server{mux: r}
}

func (s *Server) Handler() stdhttp.Handler { return s.mux }
