package events

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

const (
	TypeMoneySpent               string = "MoneySpent"
	TypeMoneyReceived            string = "MoneyReceived"
	TypeMoneyTransfered          string = "MoneyTransfered"
	TypeReimbursementReceived    string = "ReimbursementReceived"
	TypeExpectedReimbursementSet string = "ExpectedReimbursementSet"
	TypeMoneyDeposited           string = "MoneyDeposited"
	TypeMoneyWithdrawn           string = "MoneyWithdrawn"
)

type MoneySpent struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

func (e MoneySpent) Type() string {
	return TypeMoneySpent
}

func (e MoneySpent) Content() any {
	return e
}

type MoneyReceived struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

func (e MoneyReceived) Type() string {
	return TypeMoneyReceived
}

func (e MoneyReceived) Content() any {
	return e
}

type MoneyTransfered struct {
	FromAccountID uuid.UUID
	FromCurrency  values.Currency
	FromAmount    decimal.Decimal
	ToAccountID   uuid.UUID
	ToCurrency    values.Currency
	ToAmount      decimal.Decimal
	Category      string
	Description   string
}

func (e MoneyTransfered) Type() string {
	return TypeMoneyTransfered
}

func (e MoneyTransfered) Content() any {
	return e
}

type ReimbursementReceived struct {
	AccountID uuid.UUID
	From      string
	Currency  values.Currency
	Amount    decimal.Decimal
}

type ExpectedReimbursementSet struct {
	AccountID uuid.UUID
	Currency  values.Currency
	Amount    decimal.Decimal
}

func (e ExpectedReimbursementSet) Type() string {
	return TypeExpectedReimbursementSet
}

func (e ExpectedReimbursementSet) Content() any {
	return e
}

func (e ReimbursementReceived) Type() string {
	return TypeReimbursementReceived
}

func (e ReimbursementReceived) Content() any {
	return e
}

type MoneyDeposited struct {
	AccountID uuid.UUID
	Currency  values.Currency
	Amount    decimal.Decimal
}

func (e MoneyDeposited) Type() string {
	return TypeMoneyDeposited
}

func (e MoneyDeposited) Content() any {
	return e
}

type MoneyWithdrawn struct {
	AccountID uuid.UUID
	Currency  values.Currency
	Amount    decimal.Decimal
}

func (e MoneyWithdrawn) Type() string {
	return TypeMoneyWithdrawn
}

func (e MoneyWithdrawn) Content() any {
	return e
}
