// internal/http-server/dto/consent.go
package dto

import (
	"multibank/backend/internal/domain"
	"time"
)

type ConsentCreateRequest struct {
	BankID      int64               `json:"bank_id"`               // указываем, с каким банком
	ClientID    string              `json:"client_id"`             // с фронта
	Permissions []domain.Permission `json:"permissions,omitempty"` // опционально, иначе дефолт
}

type ConsentResponse struct {
	ID                 int64                `json:"id"`
	RequestID          string               `json:"request_id"`
	ConsentID          *string              `json:"consent_id,omitempty"`
	Status             domain.ConsentStatus `json:"status"`
	AutoApproved       *bool                `json:"auto_approved,omitempty"`
	Permissions        []domain.Permission  `json:"permissions"`
	Reason             string               `json:"reason"`
	RequestingBank     string               `json:"requesting_bank"`
	RequestingBankName string               `json:"requesting_bank_name"`

	BankStatus           *string    `json:"bank_status,omitempty"`
	BankCreationDateTime *time.Time `json:"bank_creation_datetime,omitempty"`
	BankStatusUpdateTime *time.Time `json:"bank_status_update_datetime,omitempty"`
	BankExpirationTime   *time.Time `json:"bank_expiration_datetime,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
