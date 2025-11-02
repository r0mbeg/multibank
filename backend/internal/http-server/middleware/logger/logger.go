package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

// New возвращает chi-совместимый middleware,
// который логирует каждый запрос через slog.
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	if log == nil {
		log = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		l := log.With(slog.String("component", "middleware/logger"))
		l.Info("logger middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			entry := l.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("remote", r.RemoteAddr),
				slog.String("agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
				slog.Duration("duration", duration),
			)

			switch {
			case ww.Status() >= 500:
				entry.Error("request completed")
			case ww.Status() >= 400:
				entry.Warn("request completed")
			default:
				entry.Info("request completed")
			}
		})
	}
}
