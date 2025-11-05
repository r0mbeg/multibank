package dto

import (
	"multibank/backend/internal/domain"
	"time"
)

// Данные для создания запроса на согласие.
//
// # Соответствует телу запроса
//
// POST /account-consents/request Open Banking API
type ConsentRequest struct {
	ClientId           string              `json:"id"`
	Permissions        []domain.Permission `json:"permissions"`
	Reason             string              `json:"reason"`
	RequestingBank     string              `json:"requesting_bank"`
	RequestingBankName string              `json:"requesting_bank_name"`
}

type ConsentRequestResponse struct {
	ConsentId    string               `json:"consent_id"`
	Status       domain.ConsentStatus `json:"status"`
	AutoApproved bool                 `json:"auto_approved"`
}

type ConsentView struct {
	ID                   string               `json:"consentId"`
	Status               domain.ConsentStatus `json:"status"`
	CreationDateTime     time.Time            `json:"creationDateTime"`
	StatusUpdateDateTime time.Time            `json:"statusUpdateDateTime"`
	Permissions          []domain.Permission  `json:"permissions"`
	ExpirationDateTime   time.Time            `json:"expirationDateTime"`
}

// Представляет данные о текущем состоянии согласия
//
// # Сответствует ответу на запрос
//
// GET /account-consents/{consent_id}
type ConsentViewWrapper struct {
	Data  ConsentView       `json:"data"`
	Links map[string]string `json:"links"`
	Meta  map[string]string `json:"meta"`
}
