// internal/domain/consent.go
package domain

import "time"

type ConsentStatus string

const (
	AwaitingAuthorisation ConsentStatus = "AwaitingAuthorization"
	Rejected              ConsentStatus = "Rejected"
	Authorised            ConsentStatus = "Authorized"
	Revoked               ConsentStatus = "Revoked"
)

type Permission string

const (
	ReadAccountsDetail     Permission = "ReadAccountsDetail"
	ReadBalances           Permission = "ReadBalances"
	ReadTransactionsDetail Permission = "ReadTransactionsDetail"
)

type AccountConsent struct {
	ID       int64
	UserID   int64
	BankID   int64
	BankCode string

	// from POST-response:
	RequestID    string
	ConsentID    *string
	Status       ConsentStatus
	AutoApproved *bool

	// client_id
	ClientID string

	// always the same permissions
	Permissions        []Permission
	Reason             string
	RequestingBank     string
	RequestingBankName string

	// dates for consent status
	CreationDateTime     *time.Time
	StatusUpdateDateTime *time.Time
	ExpirationDateTime   *time.Time

	// internal timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}
