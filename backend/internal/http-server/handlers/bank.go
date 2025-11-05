// internal/http-server/handlers/bank.go
package handlers

import (
	"context"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// thin interface with only 1 method
type Bank interface {
	ListEnabled(ctx context.Context) ([]domain.Bank, error)
}

type BankHandler struct {
	svc Bank
}

// RegisterBankRoutes registers Bank handlers
// JWT is attached in server.go to the /banks
func RegisterBankRoutes(r chi.Router, svc Bank) {
	h := &BankHandler{svc: svc}
	r.Get("/", h.GetBanks)
}

// GetBanks godoc
// @Summary      Get a list of available banks
// @Description  Retrieves the list of all banks enabled (is_enabled = true) in the system.
// @Tags         banks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {array}  dto.BankResponse  "List of banks"
// @Failure      403  {object}  dto.ErrorResponse "Access denied"
// @Failure      500  {object}  dto.ErrorResponse "Internal server error"
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
		out = append(out, dto.BankResponse{
			ID: b.ID, Name: b.Name, Code: b.Code, APIBaseURL: b.APIBaseURL, IsEnabled: b.IsEnabled,
		})
	}

	httputils.WriteJSON(w, http.StatusOK, out)

}
