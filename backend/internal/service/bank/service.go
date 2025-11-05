package bank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"multibank/backend/internal/domain"

	"log/slog"
)

type Service struct {
	log        *slog.Logger
	repo       Repository
	httpClient *http.Client
	expirySkew time.Duration // time reserve, so as not to hit the expiration end-to-end
}

type Repository interface {
	ListEnabledBanks(ctx context.Context) ([]domain.Bank, error)
	GetBankByID(ctx context.Context, id int64) (domain.Bank, error)

	UpsertBankToken(ctx context.Context, t domain.BankToken) error
	GetBankToken(ctx context.Context, bankID int64) (domain.BankToken, error)
}

func New(log *slog.Logger, repository Repository) *Service {
	return &Service{
		log:        log,
		repo:       repository,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		expirySkew: 2 * time.Minute,
	}
}

type bankTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ClientID    string `json:"client_id"`
	Algorithm   string `json:"algorithm"`
	ExpiresIn   int64  `json:"expires_in"` // seconds
}

func (s *Service) GetOrRefreshToken(ctx context.Context, bankID int64) (string, time.Time, error) {
	// 1) если в кэше есть ещё валидный — используем
	if cached, err := s.repo.GetBankToken(ctx, bankID); err == nil {
		if time.Now().Add(s.expirySkew).Before(cached.ExpiresAt) {
			return cached.AccessToken, cached.ExpiresAt, nil
		}
	}

	// 2) тянем банк и запрашиваем новый токен
	b, err := s.repo.GetBankByID(ctx, bankID)
	if err != nil {
		return "", time.Time{}, err
	}
	u, err := url.Parse(b.APIBaseURL)
	if err != nil {
		return "", time.Time{}, err
	}
	u.Path = "/auth/bank-token"

	q := u.Query()
	q.Set("client_id", b.Login)
	q.Set("client_secret", b.Password)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("bank-token %d: %s", resp.StatusCode, string(body))
	}

	var tr bankTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", time.Time{}, err
	}
	if tr.AccessToken == "" || tr.ExpiresIn <= 0 {
		return "", time.Time{}, fmt.Errorf("bank-token: invalid response")
	}

	expiresAt := time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	if err := s.repo.UpsertBankToken(ctx, domain.BankToken{
		BankID:      b.ID,
		AccessToken: tr.AccessToken,
		ExpiresAt:   expiresAt,
	}); err != nil {
		return "", time.Time{}, err
	}
	return tr.AccessToken, expiresAt, nil
}

func (s *Service) ListEnabled(ctx context.Context) ([]domain.Bank, error) {
	return s.repo.ListEnabledBanks(ctx)
}
