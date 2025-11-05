// internal/domain/bank.go

package domain

import "time"

type Bank struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	APIBaseURL string    `json:"api_base_url"`
	Login      string    `json:"-"` // don't give it out
	Password   string    `json:"-"` // don't give it out
	IsEnabled  bool      `json:"is_enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BankToken struct {
	BankID      int64     `json:"bank_id"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
