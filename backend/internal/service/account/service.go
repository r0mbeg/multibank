// internal/service/account/service.go
package account

import (
	"context"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/logger"
	ob "multibank/backend/internal/service/openbanking"
	"time"
)

type ConsentRepo interface {
	ListByUser(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountConsent, error)
}

type BankService interface {
	GetBankByID(ctx context.Context, id int64) (domain.Bank, error)
	GetOrRefreshToken(ctx context.Context, bankID int64) (string, time.Time, error)
}

type OBAccountsClient interface {
	ListAccounts(bank domain.Bank, clientID, bearer, consentID, requestingBank string) ([]ob.ListAccountsRespData, error)
	GetInterimAvailableBalance(bank domain.Bank, accountID, bearer, consentID, requestingBank string) (amount, currency string, err error)
}

type Service struct {
	log     *slog.Logger
	consent ConsentRepo
	banks   BankService
	client  OBAccountsClient
}

func New(log *slog.Logger, consentRepo ConsentRepo, banks BankService, client OBAccountsClient) *Service {
	return &Service{log: log, consent: consentRepo, banks: banks, client: client}
}

// ListUserAccounts collects accounts based on the user's consent.
// If bankID != nil filter by 1 bank
func (s *Service) ListUserAccounts(ctx context.Context, userID int64, bankID *int64) ([]domain.AccountShort, error) {
	// 1) Берём согласия пользователя (лучше фильтровать по статусу Authorized, но можно и все)
	consents, err := s.consent.ListByUser(ctx, userID, bankID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AccountShort, 0, 16)

	for _, c := range consents {
		// check consent nil
		if c.ConsentID == nil || *c.ConsentID == "" {
			continue
		}

		bank, err := s.banks.GetBankByID(ctx, c.BankID)
		if err != nil {
			s.log.Warn("get bank failed", logger.Err(err), slog.Int64("bank_id", c.BankID))
			continue
		}
		token, _, err := s.banks.GetOrRefreshToken(ctx, bank.ID)
		if err != nil {
			s.log.Warn("get token failed", logger.Err(err), slog.Int64("bank_id", c.BankID))
			continue
		}

		// 2) list of accounts by client_id + consent headers
		accs, err := s.client.ListAccounts(bank, c.ClientID, token, *c.ConsentID, c.RequestingBank)
		if err != nil {
			s.log.Warn("list accounts failed", logger.Err(err), slog.Int64("bank_id", c.BankID))
			continue
		}

		// 3) the balance is InterimAvailable (упрощение)
		for _, a := range accs {
			amount, currency, err := s.client.GetInterimAvailableBalance(bank, a.AccountID, token, *c.ConsentID, c.RequestingBank)
			if err != nil {
				s.log.Warn("get balance failed",
					logger.Err(err),
					slog.String("account_id", a.AccountID),
					slog.Int64("bank_id", c.BankID),
				)
				// skip - do not add blank balance
			}

			out = append(out, domain.AccountShort{
				AccountID:      a.AccountID,
				Nickname:       a.Nickname,
				Status:         a.Status,
				AccountSubType: a.AccountSubType,
				OpeningDate:    a.OpeningDate,
				Amount:         amount,
				Currency:       currency,
				BankCode:       c.BankCode,
				ClientID:       c.ClientID,
			})
		}
	}
	return out, nil
}
