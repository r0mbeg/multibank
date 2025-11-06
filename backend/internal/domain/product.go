// internal/domain/product.go
package domain

import "time"

type ProductType string

const (
	ProductDeposit ProductType = "deposit"
	ProductLoan    ProductType = "loan"
	ProductCard    ProductType = "card"
	ProductAccount ProductType = "account"
)

type Product struct {
	// from bank
	ProductID    string      `json:"productId"`
	ProductType  ProductType `json:"productType"`
	ProductName  string      `json:"productName"`
	Description  string      `json:"description,omitempty"`
	InterestRate float64     `json:"interestRate,omitempty"`
	MinAmount    float64     `json:"minAmount,omitempty"`
	MaxAmount    float64     `json:"maxAmount,omitempty"`
	TermMonths   int         `json:"termMonths,omitempty"`

	// meta data
	BankID    int64     `json:"bank_id"`
	BankCode  string    `json:"bank_code"`
	BankName  string    `json:"bank_name"`
	FetchedAt time.Time `json:"fetched_at"`
}

type ProductFilter struct {
	ProductType string
	BankIDs     []int64
}
