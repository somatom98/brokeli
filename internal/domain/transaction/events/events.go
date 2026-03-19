package events

import (
	"time"

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
)

type MoneySpent struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
	HappenedAt  time.Time
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
	HappenedAt  time.Time
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
	HappenedAt    time.Time
}

func (e MoneyTransfered) Type() string {
	return TypeMoneyTransfered
}

func (e MoneyTransfered) Content() any {
	return e
}

type ReimbursementReceived struct {
	AccountID   uuid.UUID
	From        string
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
	HappenedAt  time.Time
}

type ExpectedReimbursementSet struct {
	AccountID  uuid.UUID
	Currency   values.Currency
	Amount     decimal.Decimal
	HappenedAt time.Time
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
