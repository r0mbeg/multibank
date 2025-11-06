// internal/http-server/dto/bank.go

package dto

import "time"

type BankResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	APIBaseURL   string    `json:"api_base_url"`
	IsEnabled    bool      `json:"is_enabled"`
	Authorized   bool      `json:"authorized"`              // from domain.BankToken
	TokenExpires time.Time `json:"token_expires,omitempty"` // from domain.BankToken
}

type BankAuthorizeResponse struct {
	Status       string    `json:"status" example:"authorized"`
	TokenExpires time.Time `json:"token_expires"`
}
