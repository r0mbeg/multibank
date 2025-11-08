// internal/http-server/dto/consent.go
package dto

import (
	"multibank/backend/internal/domain"
	"time"
)

type ConsentCreateRequest struct {
	BankCode string `json:"bank_code"` // instead of bank_id
	ClientID string `json:"client_id"` // e.g. team014-1
	//Permissions []domain.Permission `json:"permissions,omitempty"` // always default
}

type ConsentResponse struct {
	ID                 int64                `json:"id"`
	BankCode           string               `json:"bank_code"`
	ClientID           string               `json:"client_id"`
	RequestID          string               `json:"request_id"`
	ConsentID          *string              `json:"consent_id,omitempty"`
	Status             domain.ConsentStatus `json:"status"`
	AutoApproved       *bool                `json:"auto_approved,omitempty"`
	Permissions        []domain.Permission  `json:"permissions"`
	Reason             string               `json:"reason"`
	RequestingBank     string               `json:"requesting_bank"`
	RequestingBankName string               `json:"requesting_bank_name"`

	CreationDateTime     *time.Time `json:"creation_datetime,omitempty"`
	StatusUpdateDateTime *time.Time `json:"status_update_datetime,omitempty"`
	ExpirationDateTime   *time.Time `json:"expiration_datetime,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
