package values

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Entry struct {
	AccountID uuid.UUID
	Currency  Currency
	Amount    decimal.Decimal
	Side      Side
}
