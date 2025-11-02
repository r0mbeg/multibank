package httpserver

import (
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"multibank/backend/internal/http-server/handlers"
	userSvc "multibank/backend/internal/service/user"
)

type Server struct {
	mux *chi.Mux
}

type Deps struct {
	UserService *userSvc.Service
	// Добавишь сюда другие сервисы по мере появления.
}

type Options struct {
	RequestTimeout time.Duration
}

func New(deps Deps, opts Options) *Server {
	r := chi.NewRouter()

	// базовые middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// таймаут запросов — один раз на уровне HTTP
	if opts.RequestTimeout > 0 {
		r.Use(middleware.Timeout(opts.RequestTimeout))
	}

	// Регистрация feature-роутов (каждая фича сама вешает свои пути)
	handlers.RegisterUserRoutes(r, deps.UserService)

	return &Server{mux: r}
}

func (s *Server) Handler() stdhttp.Handler { return s.mux }
