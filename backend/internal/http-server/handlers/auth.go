package handlers

import (
	"encoding/json"
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

func RegisterAuthRoutes(r chi.Router, d AuthDeps) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Email, FirstName, LastName, Patronymic, BirthDate, Password string
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json")
				return
			}
			u, err := d.Auth.Register(r.Context(), auth.RegisterInput{
				Email: req.Email, FirstName: req.FirstName, LastName: req.LastName,
				Patronymic: req.Patronymic, BirthDate: req.BirthDate, Password: req.Password,
			})
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			// Можно сразу выдать токен после регистрации
			tok, exp, err := d.JWT.Issue(u.ID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "token issue error")
				return
			}
			writeJSON(w, http.StatusCreated, map[string]any{
				"access_token": tok,
				"expires_in":   int64(time.Until(exp).Seconds()),
			})
		})

		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			var req struct{ Email, Password string }
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json")
				return
			}
			tok, err := d.Auth.Login(r.Context(), req.Email, req.Password)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid credentials")
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"access_token": tok})
		})
	})
}
