package dto

type AccountResponse struct {
	AccountID      string `json:"account_id"`
	Nickname       string `json:"nickname"`
	Status         string `json:"status"`
	AccountSubType string `json:"account_sub_type"`
	OpeningDate    string `json:"opening_date"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	BankCode       string `json:"bank_code"`
	ClientID       string `json:"client_id"`
}
