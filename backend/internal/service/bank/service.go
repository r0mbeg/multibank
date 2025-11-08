// internal/service/bank/service.go

package bank

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/logger"
	"multibank/backend/internal/storage"
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
	GetBankByCode(ctx context.Context, code string) (domain.Bank, error)

	UpsertBankToken(ctx context.Context, t domain.BankToken) error
	GetBankToken(ctx context.Context, bankID int64) (domain.BankToken, error)
}

var (
	ErrBanksNotFound = errors.New("banks not found")
)

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

	const op = "service.bank.GetOrRefreshToken"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("bank_id", bankID),
	)

	log.Info("fetching bank token")

	// 1) if there is still a valid token in the DB, we use it
	if cached, err := s.repo.GetBankToken(ctx, bankID); err == nil {
		if time.Now().Add(s.expirySkew).Before(cached.ExpiresAt) {
			return cached.AccessToken, cached.ExpiresAt, nil
		}
	}

	// 2) тянем банк и запрашиваем новый токен
	b, err := s.repo.GetBankByID(ctx, bankID)
	if err != nil {
		log.Warn("failed to get bank details", logger.Err(err))
		return "", time.Time{}, err
	}
	base, err := httputils.NormalizeURL(b.APIBaseURL)
	if err != nil {
		log.Warn("invalid bank api_base_url", slog.String("api_base_url", b.APIBaseURL), logger.Err(err))
		return "", time.Time{}, err
	}

	// /auth/bank-token + query ?client_id=...&client_secret=...
	tokenURL := base.ResolveReference(&url.URL{Path: "/auth/bank-token"})
	q := tokenURL.Query()
	q.Set("client_id", b.Login)
	q.Set("client_secret", b.Password)
	tokenURL.RawQuery = q.Encode()

	// log w/o secret
	masked := *tokenURL
	mq := masked.Query()
	mq.Set("client_secret", "******")
	masked.RawQuery = mq.Encode()
	log.Info("requesting bank token", slog.String("url", masked.String()))

	// POST w/o body
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL.String(), http.NoBody)
	if err != nil {
		log.Warn("unable to create http request", logger.Err(err))
		return "", time.Time{}, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Warn("unable send req for a bank token", logger.Err(err))
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Warn("token endpoint returned non-200", slog.String("status", resp.Status))
		return "", time.Time{}, fmt.Errorf("bank-token %d: %s", resp.StatusCode, string(body))
	}

	var tr bankTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		log.Warn("failed to decode bank token response", logger.Err(err))
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

	log.Info("bank token successfully refreshed")

	return tr.AccessToken, expiresAt, nil
}

// ListEnabled return a list of all banks where IsEnabled == true
func (s *Service) ListEnabled(ctx context.Context) ([]domain.Bank, error) {
	const op = "service.bank.ListEnabled"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting enabled banks")

	banks, err := s.repo.ListEnabledBanks(ctx)
	if err != nil {

		if errors.Is(err, storage.ErrBanksNotFound) {
			log.Warn("banks not found", logger.Err(err))
			return []domain.Bank{}, fmt.Errorf("%s: %w", op, ErrBanksNotFound)
		}

		log.Error("failed to get enabled banks", logger.Err(err))
		return []domain.Bank{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully got enabled banks")

	return banks, nil
}

// TokenStatus returns true/false and expiration time, without giving the token itself.
func (s *Service) TokenStatus(ctx context.Context, bankID int64) (bool, time.Time, error) {
	t, err := s.repo.GetBankToken(ctx, bankID)
	if err != nil || t.AccessToken == "" {
		return false, time.Time{}, nil
	}
	// Считаем валидным, если до экспирации остаётся хотя бы expirySkew.
	if time.Now().Add(s.expirySkew).Before(t.ExpiresAt) {
		return true, t.ExpiresAt, nil
	}
	return false, t.ExpiresAt, nil
}

// EnsureTokensForEnabled goes through all the enabled banks and gets valid token for each
func (s *Service) EnsureTokensForEnabled(ctx context.Context) error {

	const op = "service.EnsureTokensForEnabled"

	log := s.log.With(
		slog.String("op", op),
	)

	banks, err := s.repo.ListEnabledBanks(ctx)
	if err != nil {
		return err
	}
	for _, b := range banks {
		// try to get of update
		if _, _, err := s.GetOrRefreshToken(ctx, b.ID); err != nil {
			// log warn, but not crash the app
			log.Warn("failed to get/refresh bank token",
				slog.Int64("bank_id", b.ID),
				slog.String("code", b.Code),
				slog.String("name", b.Name),
				logger.Err(err),
			)
		}
	}
	return nil
}

// EnsureTokensForEnabledWithWorkers does the same as EnsureTokensForEnabled,
// but with worker pool
func (s *Service) EnsureTokensForEnabledWithWorkers(ctx context.Context, workers int) error {
	const op = "service.bank.EnsureTokensForEnabledWithWorkers"
	log := s.log.With(slog.String("op", op))

	if workers <= 0 {
		workers = 1
	}

	banks, err := s.repo.ListEnabledBanks(ctx)
	if err != nil {
		return err
	}
	if len(banks) == 0 {
		return nil
	}

	sem := make(chan struct{}, workers)
	errCh := make(chan error, len(banks))
	doneCh := make(chan struct{}, len(banks))

	for _, b := range banks {
		// захватываем слот
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sem <- struct{}{}:
		}

		// b for goroutine
		bank := b

		go func() {
			defer func() { <-sem }()

			// timeout for one bank
			oneCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()

			if _, _, err := s.GetOrRefreshToken(oneCtx, bank.ID); err != nil {
				log.Warn("failed to get/refresh bank token",
					slog.Int64("bank_id", bank.ID),
					slog.String("code", bank.Code),
					slog.String("name", bank.Name),
					logger.Err(err),
				)
				errCh <- err
			} else {
				doneCh <- struct{}{}
			}
		}()
	}

	// waiting for all
	var hadErr bool
	for i := 0; i < len(banks); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-doneCh:
			// ok
		case <-errCh:
			hadErr = true
		}
	}

	if hadErr {
		// return generalized error w/o crashing
		return fmt.Errorf("%s: some tokens failed to refresh", op)
	}
	return nil
}

// GetBankDetails returns bank row by id (thin wrapper over repo).
func (s *Service) GetBankByID(ctx context.Context, id int64) (domain.Bank, error) {
	return s.repo.GetBankByID(ctx, id)
}

func (s *Service) GetBankByCode(ctx context.Context, code string) (domain.Bank, error) {
	return s.repo.GetBankByCode(ctx, code)
}
