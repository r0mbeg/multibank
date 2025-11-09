// internal/http-server/handlers/products.go
package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
)

type Product interface {
	List(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error)
}

type ProductHandler struct {
	svc Product
}

func RegisterProductRoutes(r chi.Router, svc Product) {
	h := &ProductHandler{svc: svc}
	r.Get("/", h.List)
}

// List godoc
// @Summary      Get aggregated products
// @Description  Returns a merged list of products from all enabled banks with optional filters.
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        product_type  query   string  false  "Type: deposit|loan|card|account"
// @Param        bank_id       query   []int   false  "Repeatable bank id filter" collectionFormat=multi
// @Success      200  {array}  dto.ProductResponse
// @Failure      403  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	_, ok := authmw.UserIDFromContext(r.Context())
	if !ok {
		httputils.WriteError(w, http.StatusForbidden, "access denied")
		return
	}

	f := domain.ProductFilter{}

	f.ProductType = r.URL.Query().Get("product_type")
	if vals, ok := r.URL.Query()["bank_id"]; ok {
		for _, v := range vals {
			if v == "" {
				continue
			}
			for _, part := range strings.Split(v, ",") { // поддержим и comma-separated
				if id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64); err == nil && id > 0 {
					f.BankIDs = append(f.BankIDs, id)
				}
			}
		}
	}

	items, err := h.svc.List(r.Context(), f)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "internal")
		return
	}

	out := make([]dto.ProductResponse, 0, len(items))
	for _, p := range items {
		out = append(out, dto.ProductResponse{
			ProductID:     p.ProductID,
			ProductType:   p.ProductType,
			ProductName:   p.ProductName,
			Description:   p.Description,
			InterestRate:  p.InterestRate,
			MinAmount:     p.MinAmount,
			MaxAmount:     p.MaxAmount,
			TermMonths:    p.TermMonths,
			BankID:        p.BankID,
			BankCode:      p.BankCode,
			BankName:      p.BankName,
			FetchedAt:     p.FetchedAt,
			IsRecommended: p.IsRecommended,
		})
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}
