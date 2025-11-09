// internal/http-server/dto/recommended.go
package dto

import "time"

type RecommendedRule struct {
	ProductID   string    `json:"product_id"`
	BankCode    string    `json:"bank_code"`
	ProductType string    `json:"product_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type RecommendedUpsertRequest struct {
	ProductID   string `json:"product_id" example:"prod-abank-card-001"`
	BankCode    string `json:"bank_code" example:"abank"`
	ProductType string `json:"product_type" example:"credit_card"`
}
