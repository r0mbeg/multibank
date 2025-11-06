package domain

import "time"

type Account struct {
	ID             string                  `json:"accountId"`
	Status         AccountStatus           `json:"status"`
	Currency       string                  `json:"currency"`
	AccountType    AccountType             `json:"accountType"`
	AccountSubType string                  `json:"accountSubType"` /* не стандартизированное поле */
	Description    string                  `json:"description"`
	Nickname       string                  `json:"nickname"`    /* не стандартизированное поле */
	OpeningDate    time.Time               `json:"openingDate"` /* не стандартизированное поле */
	AccountDetails []AccountIdentification `json:"account"`
}

type AccountType string

const (
	Personal AccountType = "Personal"
	Business AccountType = "Business"
)

type AccountStatus string

const (
	Enabled  AccountStatus = "Enabled"
	Disabled AccountStatus = "Disabled"
	Deleted  AccountStatus = "Deleted"
)

type AccountIdentification struct {
	SchemeName     AccountScheme `json:"schemeName"`
	Identification string        `json:"identification"`
	Name           string        `json:"name"`
}

type AccountScheme string

const (
	BBAN AccountScheme = "RU.CBR.BBAN"
	EPID AccountScheme = "RU.CBR.EPID"
	PAN  AccountScheme = "RU.CBR.PAN"
	MTEL AccountScheme = "RU.CBR.MTEL"
	ORID AccountScheme = "RU.CBR.ORID"
)

type AccountBalance struct {
	AccountID   string                `json:"accountId"`
	Type        AccountBalanceType    `json:"type"`
	DateTime    time.Time             `json:"dateTime"`
	Amount      AmountDetails         `json:"amount"`
	CreditDebit CreditDebitIndicatior `json:"creditDebitIndicator"`
}

type AccountBalanceType string

const (
	InterimAvailable AccountBalanceType = "InterimAvailable"
	InterimBooked    AccountBalanceType = "InterimBooked"
	// TODO есть еще, но мне впадлу
)

type AmountDetails struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CreditDebitIndicatior string

const (
	Credit CreditDebitIndicatior = "Credit"
	Debit  CreditDebitIndicatior = "Debit"
)

type Transaction struct {
	ID                     string                `json:"transactionId"`
	AccountID              string                `json:"accountId"`
	Amount                 AmountDetails         `json:"amount"`
	CreditDebit            CreditDebitIndicatior `json:"creditDebitIndicator"`
	Status                 string                `json:"status"`
	BookingDateTime        time.Time             `json:"bookingDateTime"`
	ValueDateTime          time.Time             `json:"valueDateTime"`
	TransactionInformation string                `json:"transactionInformation"`
	BankTransactionCode    BankTransactionCode   `json:"bankTransactionCode"`
}

type BankTransactionCode struct {
	Code string `json:"code"`
}
