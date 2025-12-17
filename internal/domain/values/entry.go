package values

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Entry struct {
	AccountID uuid.UUID
	Currency  Currency
	Amount    decimal.Decimal
	Side      Side
}
