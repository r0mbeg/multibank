package handlers

import (
	"context"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Account interface {
	ListUserAccounts(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountShort, error)
}

type AccountsHandler struct {
	svc Account
}

func RegisterAccountRoutes(r chi.Router, svc Account) {
	h := &AccountsHandler{svc: svc}
	r.Get("/", h.list) // /accounts?bank_id=...
}

// list returns a list of user accounts aggregated from connected banks.
// @Summary      List user accounts
// @Description  Returns the list of accounts available for the authorized user, aggregated across all connected banks.
// @Description  Each account includes nickname, status, subtype, opening date, and current balance (InterimAvailable).
// @Tags         Accounts
// @Security     BearerAuth
// @Produce      json
// @Param        bank_id  query     int64  false  "Filter by bank ID (optional)"
// @Success      200      {array}   dto.AccountResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /accounts [get]
func (h *AccountsHandler) list(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmw.UserIDFromContext(r.Context())
	if !ok {
		httputils.WriteError(w, http.StatusUnauthorized, "missing user in context")
		return
	}

	var bankID *int64
	if s := r.URL.Query().Get("bank_id"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			bankID = &v
		}
	}

	items, err := h.svc.ListUserAccounts(r.Context(), userID, bankID)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]dto.AccountResponse, 0, len(items))
	for _, it := range items {
		out = append(out, dto.AccountResponse{
			AccountID:      it.AccountID,
			Nickname:       it.Nickname,
			Status:         it.Status,
			AccountSubType: it.AccountSubType,
			OpeningDate:    it.OpeningDate,
			Amount:         it.Amount,
			Currency:       it.Currency,
			BankCode:       it.BankCode,
			ClientID:       it.ClientID,
		})
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}
