package account

import (
	"fmt"
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
			event, err := event_store.DecodeEvent[events.Created](record.Content())
			if err != nil {
				return fmt.Errorf("decode Created event: %w", err)
			}
			a.ApplyCreated(event)
		case events.Type_MoneyDeposited:
			event, err := event_store.DecodeEvent[events.MoneyDeposited](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyDeposited event: %w", err)
			}
			a.ApplyMoneyDeposited(event)
		case events.Type_AccountClosed:
			event, err := event_store.DecodeEvent[events.AccountClosed](record.Content())
			if err != nil {
				return fmt.Errorf("decode AccountClosed event: %w", err)
			}
			a.ApplyAccountClosed(event)
		}
	}

	return nil
}
