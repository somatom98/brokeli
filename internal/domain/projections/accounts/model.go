package accounts

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type Account struct {
	Balance                map[values.Currency]decimal.Decimal `json:"balance"`
	ExpectedReimbursements map[values.Currency]decimal.Decimal `json:"expected_reimbursements"`
	CreatedAt              *time.Time                          `json:"created_at,omitempty"`
	ClosedAt               *time.Time                          `json:"closed_at,omitempty"`
}

func NewAccount() Account {
	return Account{
		Balance:                make(map[values.Currency]decimal.Decimal),
		ExpectedReimbursements: make(map[values.Currency]decimal.Decimal),
	}
}
