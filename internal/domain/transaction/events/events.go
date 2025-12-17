package events

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

const (
	Type_MoneySpent               string = "MoneySpent"
	Type_MoneyReceived            string = "MoneyReceived"
	Type_MoneyTransfered          string = "MoneyTransfered"
	Type_ReimbursementReceived    string = "ReimbursementReceived"
	Type_ExpectedReimbursementSet string = "ExpectedReimbursementSet"
)

type MoneySpent struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

func (e MoneySpent) Type() string {
	return Type_MoneySpent
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
	return Type_MoneyReceived
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
	return Type_MoneyTransfered
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
	return Type_ExpectedReimbursementSet
}

func (e ExpectedReimbursementSet) Content() any {
	return e
}

func (e ReimbursementReceived) Type() string {
	return Type_ReimbursementReceived
}

func (e ReimbursementReceived) Content() any {
	return e
}
