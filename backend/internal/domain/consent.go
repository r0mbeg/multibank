// internal/domain/consent.go

package domain

import "time"

type ConsentStatus string

const (
	AwaitingAuthorisation ConsentStatus = "AwaitingAuthorisation"
	Rejected              ConsentStatus = "Rejected"
	Authorised            ConsentStatus = "Authorised"
	Revoked               ConsentStatus = "Revoked"
)

type Permission string

const (
	ReadAccountsDetail     Permission = "ReadAccountsDetail"
	ReadBalances           Permission = "ReadBalances"
	ReadTransactionsDetail Permission = "ReadTransactionsDetail"
)

type AccountConsent struct {
	ID                 int64 // internal PK
	UserID             int64
	BankID             int64
	RequestID          string        // req-... (returns immediately)
	ConsentID          *string       // consent-... (may be got later)
	Status             ConsentStatus // AwaitingAuthorisation/Authorised and etc
	AutoApproved       *bool
	Permissions        []Permission // JSON in SQLite
	Reason             string
	RequestingBank     string
	RequestingBankName string

	CreatedAt time.Time
	UpdatedAt time.Time
	// last known from bank (raw, for debugging/tracing)
	BankStatus           *string // e.g. "AwaitingAuthorization", "Authorized"
	BankCreationDateTime *time.Time
	BankStatusUpdateTime *time.Time
	BankExpirationTime   *time.Time
}
