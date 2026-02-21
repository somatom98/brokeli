package values

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Entry struct {
	AccountID uuid.UUID
	Currency  Currency
	Amount    decimal.Decimal
	Side      Side
}

func (e Entry) String() string {
	return fmt.Sprintf("account: %s, currency: %s, amount: %s, side: %v", e.AccountID, e.Currency, e.Amount, e.Side)
}
