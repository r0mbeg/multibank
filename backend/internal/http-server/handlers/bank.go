package handlers

import (
	"context"
	"encoding/json"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	"multibank/backend/internal/service/bank"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// thin interface with only 1 method
type Bank interface {
	ListEnabled(ctx context.Context) ([]domain.Bank, error)
}

type BankHandler struct {
	svc bank.Service
}

// RegisterBankRoutes registers Bank handlers
// JWT is attached in server.go to the /banks
func RegisterBankRoutes(r chi.Router, svc bank.Service) {
	h := &BankHandler{svc: svc}
	r.Get("/", h.GetBanks)
}

func (h *BankHandler) GetBanks(w http.ResponseWriter, r *http.Request) {
	banks, err := h.svc.ListEnabled(r.Context())
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	out := make([]dto.BankResponse, 0, len(banks))
	for _, b := range banks {
		out = append(out, dto.BankResponse{
			ID: b.ID, Name: b.Name, Code: b.Code, APIBaseURL: b.APIBaseURL, IsEnabled: b.IsEnabled,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
