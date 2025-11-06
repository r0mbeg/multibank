package openbanking

import (
	"log/slog"
	"multibank/backend/internal/domain"
	"net/http"
)

type AccountOperations interface {
	GetAllAccounts(
		AccountConsent *domain.FullAccountConsent,
		RequestingBank *domain.Bank,
	) ([]*domain.Account, error)

	GetAccountByID(
		AccountConsent *domain.FullAccountConsent,
		RequestingBank *domain.Bank,
		AccountID string,
	) (*domain.Account, error)

	GetAccountBalanceByID(
		AccountConsent domain.FullAccountConsent,
		RequestingBank *domain.Bank,
		AccountID string,
	) (*domain.AccountBalance, error)

	GetAccountTransactionsByID(
		AccountConsent domain.FullAccountConsent,
		RequestingBank *domain.Bank,
		AccountID string,
	) ([]*domain.Transaction, error)
}

type AccountClient struct {
	httpClient *http.Client
	log        *slog.Logger
}
