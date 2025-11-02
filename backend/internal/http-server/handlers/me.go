package handlers

import (
	authmw "multibank/backend/internal/auth/middleware"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/service/user"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MeHandler struct {
	svc *user.Service
}

// RegisterMeRoutes registers ME handlers
// JWT is attached in server.go to the /me
func RegisterMeRoutes(r chi.Router, svc *user.Service) {
	h := MeHandler{svc: svc}
	r.Get("/", h.GetMe)
}

// GetMe godoc
// @Summary      Get current user
// @Description  Возвращает профиль текущего аутентифицированного пользователя
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.User
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /me [get]
func (h *MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	userID, ok := authmw.UserIDFromContext(r.Context())

	if !ok {
		httputils.WriteError(w, http.StatusForbidden, "access denied")
		return
	}

	u, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, dto.UserFromDomain(u))
}
