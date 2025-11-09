// internal/http-server/handlers/recommended.go
package handlers

import (
	"context"
	"encoding/json"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/service/product"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Сервисный интерфейс
type Recommended interface {
	List(ctx context.Context) ([]product.Rule, error)
	Upsert(ctx context.Context, productID, bankCode, productType string) error
	Delete(ctx context.Context, productID, bankCode, productType string) error
}

type RecommendedHandler struct {
	svc Recommended
}

func RegisterRecommendedRoutes(r chi.Router, svc Recommended) {
	h := &RecommendedHandler{svc: svc}

	// admin endpoints (повесь свою аутентификацию/авторизацию)
	r.Get("/", h.list)
	r.Post("/", h.upsert)
	r.Delete("/", h.delete)
}

// @Summary      List recommended product rules
// @Tags         admin/products
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  dto.RecommendedRule
// @Router       /admin/recommended-products [get]
func (h *RecommendedHandler) list(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.List(r.Context())
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]dto.RecommendedRule, 0, len(rows))
	for _, v := range rows {
		out = append(out, dto.RecommendedRule{
			ProductID:   v.ProductID,
			BankCode:    v.BankCode,
			ProductType: v.ProductType,
			CreatedAt:   v.CreatedAt,
		})
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}

// @Summary      Upsert recommended rule
// @Tags         admin/products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.RecommendedUpsertRequest  true  "Rule triplet"
// @Success      204    "No Content"
// @Failure      400    {object} dto.ErrorResponse
// @Failure      500    {object} dto.ErrorResponse
// @Router       /admin/recommended-products [post]
func (h *RecommendedHandler) upsert(w http.ResponseWriter, r *http.Request) {
	var req dto.RecommendedUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "bad json: "+err.Error())
		return
	}
	if req.ProductID == "" || req.BankCode == "" || req.ProductType == "" {
		httputils.WriteError(w, http.StatusBadRequest, "product_id, bank_code, product_type are required")
		return
	}
	if err := h.svc.Upsert(r.Context(), req.ProductID, req.BankCode, req.ProductType); err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Delete recommended rule
// @Tags         admin/products
// @Security     BearerAuth
// @Accept       json
// @Param        input  body  dto.RecommendedUpsertRequest  true  "Rule triplet to delete"
// @Success      204    "No Content"
// @Failure      400    {object} dto.ErrorResponse
// @Failure      500    {object} dto.ErrorResponse
// @Router       /admin/recommended-products [delete]
func (h *RecommendedHandler) delete(w http.ResponseWriter, r *http.Request) {
	var req dto.RecommendedUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "bad json: "+err.Error())
		return
	}
	if req.ProductID == "" || req.BankCode == "" || req.ProductType == "" {
		httputils.WriteError(w, http.StatusBadRequest, "product_id, bank_code, product_type are required")
		return
	}
	if err := h.svc.Delete(r.Context(), req.ProductID, req.BankCode, req.ProductType); err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
