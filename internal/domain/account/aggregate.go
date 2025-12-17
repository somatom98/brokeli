package account

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

const Version = 0

type State int

const (
	State_Unknown State = iota
	State_Created
	State_Closed
)

type Account struct {
	ID        uuid.UUID
	State     State
	Balances  map[values.Currency]decimal.Decimal
	CreatedAt time.Time
	ClosedAt  *time.Time
}

func New(id uuid.UUID) *Account {
	return &Account{
		ID:       id,
		Balances: make(map[values.Currency]decimal.Decimal),
	}
}

func (a *Account) Hydrate(records []event_store.Record) error {
	for _, record := range records {
		switch record.Type() {
		case events.Type_Created:
			a.ApplyCreated(record.Content().(events.Created))
		case events.Type_MoneyDeposited:
			a.ApplyMoneyDeposited(record.Content().(events.MoneyDeposited))
		case events.Type_AccountClosed:
			a.ApplyAccountClosed(record.Content().(events.AccountClosed))
		}
	}

	return nil
}
