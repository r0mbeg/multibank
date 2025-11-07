package service

import (
	"context"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/http-server/dto"
)

type AccountService struct {
	log           *slog.Logger
	accountClient AccountOperations
	consents      Consents
	ctx           *context.Context
}

type Consents interface {
	GetOrRefreshConsent(ctx context.Context, bankID int64, userID int64) (*domain.FullAccountConsent, error)
}

func (s *AccountService) GetAllAccountsWithBalances(
	user domain.User,
	bank domain.Bank,
) ([]*dto.AccountWithBalance, error) {
	consent, err := s.consents.GetOrRefreshConsent(*s.ctx, user.ID, bank.ID)
	if err != nil {
		return nil, err
	}
	accounts, err := s.accountClient.GetAllAccounts(consent, domain.ThisBank)
	if err != nil {
		return nil, err
	}
	accountsWithBalances := make([]*dto.AccountWithBalance, len(accounts))
	for _, account := range accounts {
		accountsWithBalance := s.GetBalanceByAccount(consent, account)
		accountsWithBalances = append(accountsWithBalances, accountsWithBalance)
	}
	return accountsWithBalances, nil
}

func (s *AccountService) GetBalanceByAccount(consent *domain.FullAccountConsent, account *domain.Account) *dto.AccountWithBalance {
	balance, err := s.accountClient.GetAccountBalanceByID(consent, domain.ThisBank, account.ID)
	if err != nil {
		return &dto.AccountWithBalance{
			Account: account,
			Balance: nil,
		}
	}
	return &dto.AccountWithBalance{
		Account: account,
		Balance: balance,
	}
}

func (s *AccountService) GetAccountByID(
	user domain.User,
	bank domain.Bank,
	account_id string,
) (*dto.AccountWithBalance, error) {
	consent, err := s.consents.GetOrRefreshConsent(*s.ctx, user.ID, bank.ID)
	if err != nil {
		return nil, err
	}
	account, err := s.accountClient.GetAccountByID(
		consent, domain.ThisBank, account_id,
	)
	if err != nil {
		return nil, err
	}
	balance := s.GetBalanceByAccount(consent, account)
	return balance, nil
}

func (s *AccountService) GetTransactionsByAccountID(
	user domain.User,
	bank domain.Bank,
	account_id string,
) ([]*domain.Transaction, error) {
	consent, err := s.consents.GetOrRefreshConsent(*s.ctx, user.ID, bank.ID)
	if err != nil {
		return nil, err
	}
	transactions, err := s.accountClient.GetAccountTransactionsByID(
		consent, domain.ThisBank, account_id,
	)
	if err != nil {
		return nil, err
	}
	return transactions
}
