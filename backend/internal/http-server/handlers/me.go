// internal/http-server/handlers/me.go

package handlers

import (
	"errors"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
	usersvc "multibank/backend/internal/service/user"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MeHandler struct {
	svc User
}

// RegisterMeRoutes registers ME handlers
// JWT is attached in server.go to the /me
func RegisterMeRoutes(r chi.Router, svc User) {
	h := &MeHandler{svc: svc}
	r.Get("/", h.GetMe)
}

// GetMe godoc
// @Summary      Get current user
// @Description  Возвращает профиль текущего аутентифицированного пользователя
// @Tags         me
// @Produce      json
// @Param        Authorization header string true "Bearer {token}" default(Bearer eyJhbGciOi...)
// @Security     BearerAuth
// @Success      200  {object} dto.UserResponse
// @Failure      403  {object} dto.ErrorResponse
// @Failure      404  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /me [get]

func (h *MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	userID, ok := authmw.UserIDFromContext(r.Context())

	if !ok {
		httputils.WriteError(w, http.StatusForbidden, "access denied")
		return
	}

	u, err := h.svc.GetByID(r.Context(), userID)
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
