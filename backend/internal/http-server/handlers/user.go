// internal/http-server/handlers/user.go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"multibank/backend/internal/http-server/dto"
	userSvc "multibank/backend/internal/service/user"
)

func RegisterUserRoutes(r chi.Router, svc *userSvc.Service) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", getUserByID(svc))
		// r.Post("/", createUser(svc))
		// r.Put("/{id}/names", updateNames(svc))
		// r.Delete("/{id}", deleteUser(svc))
		// r.Get("/", listUsers(svc))
	})
}

func getUserByID(svc *userSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		u, err := svc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, userSvc.ErrUserNotFound) {
				writeError(w, http.StatusNotFound, "user not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Маппим domain → DTO (при желании можно отдавать domain как есть)
		resp := dto.UserFromDomain(u)
		writeJSON(w, http.StatusOK, resp)
	}
}

// Локальные утилиты хендлеров.
// Если понадобятся в нескольких файлах — вынеси в internal/http-server/httputil.go.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
