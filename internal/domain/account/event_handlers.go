package account

import (
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func (a *Account) ApplyCreated(e events.Created) {
	a.State = State_Created
	a.CreatedAt = e.Time
}

func (a *Account) ApplyMoneyDeposited(e events.MoneyDeposited) {
	a.State = State_Created
	a.ensureBalanceExists(e.Currency)
	a.Balances[e.Currency] = a.Balances[e.Currency].Add(e.Amount)
}

func (a *Account) ApplyAccountClosed(e events.AccountClosed) {
	a.State = State_Closed
	a.ClosedAt = &e.Time
}

func (a *Account) ensureBalanceExists(currency values.Currency) {
	if _, ok := a.Balances[currency]; !ok {
		a.Balances[currency] = decimal.Zero
	}
}
