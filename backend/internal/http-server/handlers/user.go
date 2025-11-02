// internal/http-server/handlers/user.go
package handlers

import (
	"encoding/json"
	"errors"
	"multibank/backend/internal/http-server/dto"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	authmw "multibank/backend/internal/auth/middleware"
	"multibank/backend/internal/service/user"
)

// Регистрирует хендлеры пользователей (без привязки к JWT-менеджеру).
// Авторизацию вешаем в server.go на подроут /users.
func RegisterUserRoutes(r chi.Router, svc *user.Service) {
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		// userID из контекста положил JWT-middleware
		authID, ok := authmw.UserIDFromContext(r.Context())
		if !ok || authID != id {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}

		u, err := svc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				writeError(w, http.StatusNotFound, "user not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, http.StatusOK, dto.UserFromDomain(u))
	})
}

// Локальные утилиты хендлеров.
// Если понадобятся в нескольких файлах, можно вынести в internal/http-server/httputil.go.
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
