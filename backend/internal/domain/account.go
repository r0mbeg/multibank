package domain

// AccountShort — то, что нам нужно для UI (без хранения в БД)
type AccountShort struct {
	AccountID      string // /accounts.data.account[*].accountId
	Nickname       string // /accounts ... nickname
	Status         string // Enabled/Disabled (как в API)
	AccountSubType string // /accounts ... accountSubType
	OpeningDate    string // YYYY-MM-DD (как в API)
	Amount         string // из /accounts/{id}/balances -> InterimAvailable.amount.amount
	Currency       string // из /accounts/{id}/balances -> InterimAvailable.amount.currency

	BankCode string
	ClientID string
}
