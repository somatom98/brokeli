package transaction

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type CreateExpense struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

type CreateIncome struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

type CreateTransfer struct {
	FromAccountID uuid.UUID
	FromCurrency  values.Currency
	FromAmount    decimal.Decimal
	ToAccountID   uuid.UUID
	ToCurrency    values.Currency
	ToAmount      decimal.Decimal
	Category      string
	Description   string
}
