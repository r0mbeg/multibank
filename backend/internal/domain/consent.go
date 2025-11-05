package domain

import "time"

type Consent struct {
	ID                   string
	Status               ConsentStatus
	CreationDateTime     time.Time
	StatusUpdateDateTime time.Time
	Permissions          []Permission
	ExpirationDateTime   time.Time
}

type ConsentStatus string

// TODO узнать какие там на самом деле статусы
// взято из спек опенбанкинга
const (
	AwaitingAuthorisation ConsentStatus = "AwaitingAuthorisation"
	Rejected              ConsentStatus = "Rejected"
	Authorised            ConsentStatus = "Authorised"
	Revoked               ConsentStatus = "Revoked"
)

type Permission string

// TODO узнать какие еще разрешения можно передавать
const (
	ReadAccountsDetail     Permission = "ReadAccountsDetail"
	ReadBalances           Permission = "ReadBalances"
	ReadTransactionsDetail Permission = "ReadTransactionsDetail"
)
