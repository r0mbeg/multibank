package handlers

import (
	"context"
	"encoding/json"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
	httputils "multibank/backend/internal/http-server/utils"
	authmw "multibank/backend/internal/service/auth/middleware"
	"multibank/backend/internal/service/consent"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Consent interface {
	Request(ctx context.Context, in consent.CreateInput) (int64, error)
	Refresh(ctx context.Context, id int64) (domain.AccountConsent, error) // мы сделали (domain.AccountConsent, error), ниже приведём к dto
	Get(ctx context.Context, id int64) (domain.AccountConsent, error)
	ListMine(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountConsent, error)
	Delete(ctx context.Context, id int64) error

	RefreshStale(ctx context.Context, batchLimit, workers int) (int, error)
}

type ConsentHandler struct {
	svc Consent
}

func RegisterConsentRoutes(r chi.Router, svc Consent) {
	h := &ConsentHandler{svc: svc}

	// Base prefix is outside: r.Route("/consents", ...) in server.go
	r.Post("/request", h.request)
	r.Get("/", h.list) // ?bank_id=...
	r.Get("/{id}", h.get)
	r.Post("/{id}/refresh", h.refresh)
	r.Delete("/{id}", h.delete)
}

// request creates a consent request to the bank
// @Summary      Request account consent
// @Description  Creates a request for consent (account-consent) in the bank and saves the draft with us. If the bank auto—confirms, the consent_id will be returned immediately, otherwise the status will be Awaiting Authorization.
// @Tags         Consents
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.ConsentCreateRequest  true  "Create consent payload"
// @Success      201    {object}  dto.ConsentResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /consents/request [post]
func (h *ConsentHandler) request(w http.ResponseWriter, r *http.Request) {
	var req dto.ConsentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "bad json: "+err.Error())
		return
	}

	userID, ok := authmw.UserIDFromContext(r.Context())
	if !ok {
		httputils.WriteError(w, http.StatusUnauthorized, "missing user in context")
		return
	}

	id, err := h.svc.Request(r.Context(), consent.CreateInput{
		UserID:   userID,
		BankCode: req.BankCode,
		ClientID: req.ClientID,
		// Permissions do not take it from client. Using defaults in service
	})
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusCreated, toConsentResponse(c))
}

// list retrieves the list of current user's consents
// @Summary      List my consents
// @Description  Returns the user's consents. Can be filtered by bank.
// @Tags         Consents
// @Security     BearerAuth
// @Produce      json
// @Param        bank_id  query     int64  false  "Filter by bank id"
// @Success      200      {array}   dto.ConsentResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /consents [get]
func (h *ConsentHandler) list(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.svc.ListMine(r.Context(), userID, bankID)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]dto.ConsentResponse, 0, len(items))
	for _, it := range items {
		out = append(out, toConsentResponse(it))
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}

// get returns one consent by ID
// @Summary      Get consent by id
// @Tags         Consents
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int64  true  "Consent ID (internal)"
// @Success      200  {object}  dto.ConsentResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /consents/{id} [get]
func (h *ConsentHandler) get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	c, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, toConsentResponse(c))
}

// refresh updates the consent status from the bank.
// @Summary      Refresh consent status
// @Description  Asks the bank by request_id/consent_id and updates our status (and consent_id, if available).
// @Tags         Consents
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int64  true  "Consent ID (internal)"
// @Success      200  {object}  dto.ConsentResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /consents/{id}/refresh [post]
func (h *ConsentHandler) refresh(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	c, err := h.svc.Refresh(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, toConsentResponse(c))
}

// delete deletes the consent from us (without calling the bank).
// @Summary      Delete consent
// @Tags         Consents
// @Security     BearerAuth
// @Param        id   path  int64  true  "Consent ID (internal)"
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /consents/{id} [delete]
func (h *ConsentHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toConsentResponse(c domain.AccountConsent) dto.ConsentResponse {
	return dto.ConsentResponse{
		ID:        c.ID,
		RequestID: c.RequestID,
		ConsentID: c.ConsentID,

		BankCode:           c.BankCode, // from VIEW/JOIN
		Reason:             c.Reason,
		RequestingBank:     c.RequestingBank,
		RequestingBankName: c.RequestingBankName,

		Status:       c.Status,
		AutoApproved: c.AutoApproved, // *bool — DTO should accept *bool
		ClientID:     c.ClientID,
		Permissions:  c.Permissions,

		CreationDateTime:     c.CreationDateTime,
		StatusUpdateDateTime: c.StatusUpdateDateTime,
		ExpirationDateTime:   c.ExpirationDateTime,
		CreatedAt:            c.CreatedAt,
		UpdatedAt:            c.UpdatedAt,
	}
}
