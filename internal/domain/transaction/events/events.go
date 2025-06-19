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

type IncomeCreated struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
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
