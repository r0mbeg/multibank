// internal/http-server/handlers/auth.go
package handlers

import (
	"encoding/json"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	"net/http"
	"time"

	"multibank/backend/internal/auth"
	"multibank/backend/internal/auth/jwt"

	"github.com/go-chi/chi/v5"
)

type AuthDeps struct {
	Auth *auth.Service
	JWT  *jwt.Manager
}

type AuthHandler struct {
	auth *auth.Service
	jwt  *jwt.Manager
}

func RegisterAuthRoutes(r chi.Router, d AuthDeps) {
	h := &AuthHandler{auth: d.Auth, jwt: d.JWT}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})
}

// Register godoc
// @Summary      Register user
// @Description  User registration. Returns access_token (can be used instantly).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.RegisterRequest true "Register payload"
// @Success      201     {object} dto.TokenResponse
// @Failure      400     {object} dto.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, err := h.auth.Register(r.Context(), auth.RegisterInput{
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		BirthDate:  req.BirthDate,
		Password:   req.Password,
	})
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tok, exp, err := h.jwt.Issue(u.ID)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "token issue error")
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, dto.TokenResponse{
		AccessToken: tok,
		ExpiresIn:   int64(time.Until(exp).Seconds()),
	})
}

// Login godoc
// @Summary      Login
// @Description  Returns access_token using e-mail and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.LoginRequest true "Login payload"
// @Success      200     {object} dto.TokenResponse
// @Failure      400     {object} dto.ErrorResponse
// @Failure      401     {object} dto.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	tok, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httputils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, dto.TokenResponse{AccessToken: tok})
}
