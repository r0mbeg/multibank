// internal/service/openbanking/consent.go

package openbanking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"multibank/backend/internal/domain"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/logger"
	"net/http"
	"net/url"
	"time"
)

type ConsentClient struct {
	log  *slog.Logger
	HTTP *http.Client

	// Constants - from app.New(...)
	RequestingBank     string // "team014"
	RequestingBankName string // "Team 14 Multibank"
	Reason             string // "Account aggregation for HackAPI"
}

func NewConsentClient(log *slog.Logger, httpClient *http.Client, reqBank, reqBankName, reason string) *ConsentClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &ConsentClient{
		HTTP:               httpClient,
		log:                log,
		RequestingBank:     reqBank,
		RequestingBankName: reqBankName,
		Reason:             reason,
	}
}

func isOK(code int) bool { return code >= 200 && code < 300 }

type requestBody struct {
	ClientID           string              `json:"client_id"`
	Permissions        []domain.Permission `json:"permissions"`
	Reason             string              `json:"reason"`
	RequestingBank     string              `json:"requesting_bank"`
	RequestingBankName string              `json:"requesting_bank_name"`
}

type ConsentRequestResp struct {
	RequestID    string  `json:"request_id"`
	ConsentID    *string `json:"consent_id"` // can be blank
	Status       string  `json:"status"`     // "AwaitingAuthorisation" | "Authorised" ...
	Message      string  `json:"message"`
	CreatedAt    string  `json:"created_at"`
	AutoApproved *bool   `json:"auto_approved"`
}

type ConsentViewWrapper struct {
	Data struct {
		ConsentID            string              `json:"consentId"`
		Status               string              `json:"status"` // "Authorized" | "AwaitingAuthorization" ...
		CreationDateTime     time.Time           `json:"creationDateTime"`
		StatusUpdateDateTime time.Time           `json:"statusUpdateDateTime"`
		Permissions          []domain.Permission `json:"permissions"`
		ExpirationDateTime   time.Time           `json:"expirationDateTime"`
	} `json:"data"`
}

func (c *ConsentClient) RequestConsent(bank domain.Bank, clientID string, perms []domain.Permission, bearer string) (*ConsentRequestResp, error) {

	const op = "service.openbanking.RequestConsent"

	log := c.log.With(
		slog.String("op", op),
	)

	base, err := httputils.NormalizeURL(bank.APIBaseURL)
	if err != nil {
		log.Warn("invalid bank api_base_url", slog.String("api_base_url", bank.APIBaseURL), logger.Err(err))
		return nil, err
	}

	u, _ := url.JoinPath(base.String(), "account-consents", "request")

	body := requestBody{
		ClientID:           clientID,
		Permissions:        perms,
		Reason:             c.Reason,
		RequestingBank:     c.RequestingBank,
		RequestingBankName: c.RequestingBankName,
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", u, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("X-Requesting-Bank", c.RequestingBank)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Warn("failed to request consent", logger.Err(err))
		return nil, err
	}
	defer resp.Body.Close()
	if !isOK(resp.StatusCode) {
		all, _ := io.ReadAll(resp.Body)

		log.Warn("got non-ok status code from request",
			slog.Int("code", resp.StatusCode),
			slog.String("body", string(all)),
		)

		return nil, fmt.Errorf("consents request %d: %s", resp.StatusCode, string(all))
	}

	var out ConsentRequestResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Warn("failed to decode consent response", logger.Err(err))
		return nil, err
	}
	return &out, nil
}

func (c *ConsentClient) GetConsent(bank domain.Bank, requestOrConsentID string, xFapi string) (*ConsentViewWrapper, error) {

	const op = "service.openbanking.GetConsent"

	log := c.log.With(
		slog.String("op", op),
	)

	base, err := httputils.NormalizeURL(bank.APIBaseURL)
	if err != nil {
		log.Warn("invalid bank api_base_url", slog.String("api_base_url", bank.APIBaseURL), logger.Err(err))
		return nil, err
	}
	u, _ := url.JoinPath(base.String(), "account-consents", requestOrConsentID)

	req, _ := http.NewRequest("GET", u, nil)
	if xFapi != "" {
		req.Header.Set("x-fapi-interaction-id", xFapi)
	} else {
		c.log.Warn("x-fapi-interaction-id is blank")
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Warn("failed to get consent", logger.Err(err))
		return nil, err
	}
	defer resp.Body.Close()
	if !isOK(resp.StatusCode) {
		all, _ := io.ReadAll(resp.Body)
		log.Warn("got non-ok status code from request",
			slog.Int("code", resp.StatusCode),
			slog.String("body", string(all)),
		)
		return nil, fmt.Errorf("consents get %d: %s", resp.StatusCode, string(all))
	}
	var v ConsentViewWrapper
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Warn("failed to decode consent response", logger.Err(err))
		return nil, err
	}
	return &v, nil
}
