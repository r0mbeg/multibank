package dto

import "multibank/backend/internal/domain"

type APIWrapper[T any] struct {
	Data  T                 `json:"data"`
	Links map[string]string `json:"links"`
	Meta  map[string]string `json:"meta"`
}

type AccountWrapper struct {
	Accounts []*domain.Account `json:"account"`
}

type BalanceWrapper struct {
	Balances []*domain.AccountBalance `json:"balance"`
}
