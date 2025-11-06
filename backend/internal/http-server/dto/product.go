// internal/http-server/dto/product.go
package dto

import (
	"multibank/backend/internal/domain"
	"time"
)

type ProductResponse struct {
	ProductID    string             `json:"productId"`
	ProductType  domain.ProductType `json:"productType"`
	ProductName  string             `json:"productName"`
	Description  string             `json:"description,omitempty"`
	InterestRate float64            `json:"interestRate,omitempty"`
	MinAmount    float64            `json:"minAmount,omitempty"`
	MaxAmount    float64            `json:"maxAmount,omitempty"`
	TermMonths   int                `json:"termMonths,omitempty"`

	BankID    int64     `json:"bank_id"`
	BankCode  string    `json:"bank_code"`
	BankName  string    `json:"bank_name"`
	FetchedAt time.Time `json:"fetched_at"`
}
