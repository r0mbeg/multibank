package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

type Server struct {
	svc     UserService
	timeout time.Duration
	mux     *chi.Mux
}

func New(svc UserService, timeout time.Duration) *Server {
	r := chi.NewRouter()

	s := &Server{
		svc:     svc,
		timeout: timeout,
		mux:     r,
	}

	// маршруты
	r.Get("/users/{id}", s.getUserByID)

	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), s.timeout)
	defer cancel()

	u, err := s.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, u)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	type errResp struct {
		Error string `json:"error"`
	}
	writeJSON(w, status, errResp{Error: msg})
}
