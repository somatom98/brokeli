package values

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Movement struct {
	AccountID uuid.UUID
	Currency  Currency
	Amount    decimal.Decimal
	Side      Side
}
