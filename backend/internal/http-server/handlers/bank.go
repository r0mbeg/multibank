// internal/http-server/handlers/bank.go

package handlers

import (
	"context"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// thin interface of Bank Service
type Bank interface {
	ListEnabled(ctx context.Context) ([]domain.Bank, error)
	TokenStatus(ctx context.Context, bankID int64) (bool, time.Time, error)
	GetOrRefreshToken(ctx context.Context, bankID int64) (string, time.Time, error)
	EnsureTokensForEnabled(ctx context.Context) error
}

type BankHandler struct {
	svc Bank
}

// RegisterBankRoutes registers Bank handlers
// JWT is attached in server.go to the /banks
func RegisterBankRoutes(r chi.Router, svc Bank) {
	h := &BankHandler{svc: svc}
	r.Get("/", h.GetBanks)
	r.Post("/{id}/authorize", h.AuthorizeBank)

}

// GetBanks godoc
// @Summary      Get a list of available banks
// @Description  Retrieves all enabled banks and shows whether backend is authorized for each (token cached) and token expiry.
// @Tags         banks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.BankResponse   "List of banks"
// @Failure      403  {object}  dto.ErrorResponse  "Access denied"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /banks [get]
func (h *BankHandler) GetBanks(w http.ResponseWriter, r *http.Request) {
	_, ok := authmw.UserIDFromContext(r.Context())
	if !ok {
		httputils.WriteError(w, http.StatusForbidden, "access denied")
		return
	}

	banks, err := h.svc.ListEnabled(r.Context())
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "internal")
		return
	}

	out := make([]dto.BankResponse, 0, len(banks))
	for _, b := range banks {
		authorized, exp, _ := h.svc.TokenStatus(r.Context(), b.ID) // do not crash req if couldn't get status
		out = append(out, dto.BankResponse{
			ID:           b.ID,
			Name:         b.Name,
			Code:         b.Code,
			APIBaseURL:   b.APIBaseURL,
			IsEnabled:    b.IsEnabled,
			Authorized:   authorized,
			TokenExpires: exp,
		})
	}

	httputils.WriteJSON(w, http.StatusOK, out)
}

// AuthorizeBank godoc
// @Summary      Authorize backend in a specific bank
// @Description  Requests or refreshes access token for the bank and stores it in cache/DB.
// @Tags         banks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Bank ID"
// @Success      200  {object}  dto.BankAuthorizeResponse
// @Failure      400  {object}  dto.ErrorResponse  "Invalid id"
// @Failure      403  {object}  dto.ErrorResponse  "Access denied"
// @Failure      502  {object}  dto.ErrorResponse  "Bank auth failed"
// @Router       /banks/{id}/authorize [post]
func (h *BankHandler) AuthorizeBank(w http.ResponseWriter, r *http.Request) {
	// auth check like in GetBanks
	_, ok := authmw.UserIDFromContext(r.Context())
	if !ok {
		httputils.WriteError(w, http.StatusForbidden, "access denied")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	_, exp, err := h.svc.GetOrRefreshToken(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, http.StatusBadGateway, "bank auth failed")
		return
	}
	httputils.WriteJSON(w, http.StatusOK, dto.BankAuthorizeResponse{
		Status:       "authorized",
		TokenExpires: exp,
	})
}
