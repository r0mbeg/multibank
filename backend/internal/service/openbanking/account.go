package openbanking

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"multibank/backend/internal/domain"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/logger"
	"net/http"
	"net/url"
)

type AccountClient struct {
	log  *slog.Logger
	HTTP *http.Client
}

func NewAccountClient(log *slog.Logger, httpClient *http.Client) *AccountClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &AccountClient{log: log, HTTP: httpClient}
}

type listAccountsResp struct {
	Data struct {
		Account []struct {
			AccountID      string `json:"accountId"`
			Status         string `json:"status"` // do not need for us
			Currency       string `json:"currency"`
			AccountType    string `json:"accountType"`
			AccountSubType string `json:"accountSubType"`
			Nickname       string `json:"nickname"`
			OpeningDate    string `json:"openingDate"`
			// account[...] нам не нужен для выдачи сейчас
		} `json:"account"`
	} `json:"data"`
}

type balancesResp struct {
	Data struct {
		Balance []struct {
			AccountID string `json:"accountId"`
			Type      string `json:"type"` // look for "InterimAvailable"
			DateTime  string `json:"dateTime"`
			Amount    struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"amount"`
			CreditDebitIndicator string `json:"creditDebitIndicator"`
		} `json:"balance"`
	} `json:"data"`
}

// internal struct for not using dependencies
type ListAccountsRespData struct {
	AccountID      string
	Nickname       string
	Status         string
	AccountSubType string
	OpeningDate    string
}

// ListAccounts calls GET /accounts?client_id=... with HEADERS with Auth token and consent_id
func (c *AccountClient) ListAccounts(bank domain.Bank, clientID, bearer, consentID, requestingBank string) ([]ListAccountsRespData, error) {
	const op = "openbanking.accounts.ListAccounts"
	log := c.log.With(slog.String("op", op))

	base, err := httputils.NormalizeURL(bank.APIBaseURL)
	if err != nil {
		log.Warn("invalid bank api_base_url", slog.String("api_base_url", bank.APIBaseURL), logger.Err(err))
		return nil, err
	}

	// /accounts?client_id=...
	u, _ := url.JoinPath(base.String(), "accounts")
	uu, _ := url.Parse(u)
	q := uu.Query()
	q.Set("client_id", clientID)
	uu.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", uu.String(), nil)
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("x-consent-id", consentID)
	req.Header.Set("x-requesting-bank", requestingBank)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Warn("list accounts request failed", logger.Err(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		all, _ := io.ReadAll(resp.Body)
		log.Warn("list accounts non-ok", slog.Int("code", resp.StatusCode), slog.String("body", string(all)))
		return nil, fmt.Errorf("list accounts %d: %s", resp.StatusCode, string(all))
	}

	var v listAccountsResp
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Warn("list accounts decode failed", logger.Err(err))
		return nil, err
	}

	// Приведём к более удобной внутренней структуре
	out := make([]ListAccountsRespData, 0, len(v.Data.Account))
	for _, a := range v.Data.Account {
		out = append(out, ListAccountsRespData{
			AccountID:      a.AccountID,
			Nickname:       a.Nickname,
			Status:         a.Status,
			AccountSubType: a.AccountSubType,
			OpeningDate:    a.OpeningDate,
		})
	}
	return out, nil
}

func (c *AccountClient) GetInterimAvailableBalance(bank domain.Bank, accountID, bearer, consentID, requestingBank string) (amount, currency string, err error) {
	const op = "openbanking.accounts.GetInterimAvailableBalance"
	log := c.log.With(slog.String("op", op))

	base, err := httputils.NormalizeURL(bank.APIBaseURL)
	if err != nil {
		log.Warn("invalid bank api_base_url", slog.String("api_base_url", bank.APIBaseURL), logger.Err(err))
		return "", "", err
	}

	u, _ := url.JoinPath(base.String(), "accounts", accountID, "balances")

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("x-consent-id", consentID)
	req.Header.Set("x-requesting-bank", requestingBank)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Warn("get balances request failed", logger.Err(err))
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		all, _ := io.ReadAll(resp.Body)
		log.Warn("get balances non-ok", slog.Int("code", resp.StatusCode), slog.String("body", string(all)))
		return "", "", fmt.Errorf("get balances %d: %s", resp.StatusCode, string(all))
	}

	var v balancesResp
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Warn("get balances decode failed", logger.Err(err))
		return "", "", err
	}

	for _, b := range v.Data.Balance { // ← было v.Data.balance
		if b.Type == "InterimAvailable" {
			return b.Amount.Amount, b.Amount.Currency, nil
		}
	}
	// If no InterimAvailable return blank str
	return "", "", nil
}
