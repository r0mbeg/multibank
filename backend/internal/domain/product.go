package domain

import "time"

type Product struct {
	// from bank
	ProductID    string  `json:"productId"`
	ProductType  string  `json:"productType"` // was ProductType
	ProductName  string  `json:"productName"`
	Description  string  `json:"description,omitempty"`
	InterestRate float64 `json:"interestRate,omitempty"`
	MinAmount    float64 `json:"minAmount,omitempty"`
	MaxAmount    float64 `json:"maxAmount,omitempty"`
	TermMonths   int     `json:"termMonths,omitempty"`

	// meta data
	BankID    int64     `json:"bank_id"`
	BankCode  string    `json:"bank_code"`
	BankName  string    `json:"bank_name"`
	FetchedAt time.Time `json:"fetched_at"`

	IsRecommended bool `json:"is_recommended"`
}

type ProductFilter struct {
	ProductType string
	BankIDs     []int64
}
