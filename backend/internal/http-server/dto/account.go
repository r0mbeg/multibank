package dto

import "multibank/backend/internal/domain"

type AccountWithBalance struct {
	Account *domain.Account
	Balance *domain.AccountBalance
}
