// internal/http-server/handlers/auth.go

package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/service/auth"
	"multibank/backend/internal/service/auth/jwt"
	usrsvc "multibank/backend/internal/service/user"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Interface Auth describes what handler needs from auth layer
type Auth interface {
	Register(ctx context.Context, in auth.RegisterInput) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	auth Auth
	jwt  *jwt.Manager
}

/*
func NewAuthHandler(a Auth, j *jwt.Manager) *AuthHandler {
	return &AuthHandler{auth: a, jwt: j}
}*/

func RegisterAuthRoutes(r chi.Router, a Auth, j *jwt.Manager) {
	h := &AuthHandler{auth: a, jwt: j}
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	/*
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Register)
			r.Post("/login", h.Login)
		})*/
}

// Register godoc
// @Summary      Register user
// @Description  User registration. Returns access_token (usable immediately).
// @Description
// @Description  **Request example**
// @Description  ```json
// @Description  {
// @Description    "email": "user@example.com",
// @Description    "first_name": "Ivan",
// @Description    "last_name": "Petrov",
// @Description    "patronymic": "Ivanovich",
// @Description    "birthdate": "1990-01-15",
// @Description    "password": "P@ssw0rd123"
// @Description  }
// @Description  ```
// @Description
// @Description  **Response example**
// @Description  ```json
// @Description  {
// @Description    "access_token": "eyJhbGciOi...",
// @Description    "expires_in": 3600
// @Description  }
// @Description  ```
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.RegisterRequest true "Register payload"
// @Success      201     {object} dto.TokenResponse
// @Failure      400     {object} dto.ErrorResponse
// @Failure      409     {object} dto.ErrorResponse "email already used"
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
		switch {
		case errors.Is(err, usrsvc.ErrEmailAlreadyUsed):
			httputils.WriteError(w, http.StatusConflict, "email already used")
		default:
			httputils.WriteError(w, http.StatusBadRequest, "registration failed")
		}
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
// @Description
// @Description  **Request example**
// @Description  ```json
// @Description  {
// @Description    "email": "user@example.com",
// @Description    "password": "P@ssw0rd123"
// @Description  }
// @Description  ```
// @Description
// @Description  **Response example**
// @Description  ```json
// @Description  {
// @Description    "access_token": "eyJhbGciOi...",
// @Description    "expires_in": 3600
// @Description  }
// @Description  ```
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
		httputils.WriteError(w, http.StatusUnauthorized, auth.ErrInvalidCredentials.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, dto.TokenResponse{AccessToken: tok})
}
