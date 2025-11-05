// internal/http-server/handlers/user.go
package handlers

import (
	"context"
	"errors"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
	"net/http"
	"strconv"

	usersvc "multibank/backend/internal/service/user"

	"github.com/go-chi/chi/v5"
)

// thin interface with only 1 method
type User interface {
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

type UserHandler struct {
	svc User
}

// RegisterUserRoutes registers user handlers
// JWT is attached in server.go to the /users
func RegisterUserRoutes(r chi.Router, svc User) {
	h := &UserHandler{svc: svc}
	r.Get("/{id}", h.GetByID)
	// ещё какие-то хендлеры
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Доступ ограничён владельцем токена (запрещён доступ к чужим профилям).
// @Tags         users
// @Produce      json
// @Param        Authorization header string true "Bearer {token}" default(Bearer eyJhbGciOi...)
// @Security     BearerAuth
// @Param        id   path      int64 true "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/{id} [get]

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// userID from Context (which was placed by JWT-middleware)
	// можно смотреть только нашего пользователя
	authID, ok := authmw.UserIDFromContext(r.Context())
	if !ok || authID != id {
		httputils.WriteError(w, http.StatusForbidden, "access denied")
		return
	}

	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, usersvc.ErrUserNotFound) {
			httputils.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, dto.UserResponseFromDomain(u))

}
