package dto

type BankResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	APIBaseURL string `json:"api_base_url"`
	IsEnabled  bool   `json:"is_enabled"`
}
