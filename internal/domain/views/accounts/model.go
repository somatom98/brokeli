package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type Account struct {
	Balance map[values.Currency]decimal.Decimal `json:"balance"`
}

func NewAccount() Account {
	return Account{
		Balance: make(map[values.Currency]decimal.Decimal),
	}
}
