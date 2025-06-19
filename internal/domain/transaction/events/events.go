package events

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type ExpenseCreated struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

func (e ExpenseCreated) Type() string {
	return "ExpenseCreated"
}

func (e ExpenseCreated) Content() any {
	return e
}

type IncomeCreated struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

func (e IncomeCreated) Type() string {
	return "IncomeCreated"
}

func (e IncomeCreated) Content() any {
	return e
}

type TransferCreated struct {
	FromAccountID uuid.UUID
	FromCurrency  values.Currency
	FromAmount    decimal.Decimal
	ToAccountID   uuid.UUID
	ToCurrency    values.Currency
	ToAmount      decimal.Decimal
	Category      string
	Description   string
}

func (e TransferCreated) Type() string {
	return "TransferCreated"
}

func (e TransferCreated) Content() any {
	return e
}
