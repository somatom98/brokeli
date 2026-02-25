package account

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type State int

const (
	State_Unopened State = iota
	State_Opened
)

type Account struct {
	ID    uuid.UUID
	State State
}

func New(id uuid.UUID) *Account {
	return &Account{
		ID:    id,
		State: State_Unopened,
	}
}

func (a *Account) Hydrate(records []event_store.Record) error {
	for _, record := range records {
		switch record.Type() {
		case events.TypeOpened:
			event, err := event_store.DecodeEvent[events.Opened](record.Content())
			if err != nil {
				return fmt.Errorf("decode Opened event: %w", err)
			}
			a.ApplyOpened(event)
		case events.TypeNameUpdated:
			event, err := event_store.DecodeEvent[events.NameUpdated](record.Content())
			if err != nil {
				return fmt.Errorf("decode NameUpdated event: %w", err)
			}
			a.ApplyNameUpdated(event)
		case events.TypeMoneyDeposited:
			event, err := event_store.DecodeEvent[events.MoneyDeposited](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyDeposited event: %w", err)
			}
			a.ApplyMoneyDeposited(event)
		case events.TypeMoneyWithdrawn:
			event, err := event_store.DecodeEvent[events.MoneyWithdrawn](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyWithdrawn event: %w", err)
			}
			a.ApplyMoneyWithdrawn(event)
		}
	}

	return nil
}

func (a *Account) ApplyOpened(event events.Opened) {
	a.ID = event.AccountID
	a.State = State_Opened
}

func (a *Account) ApplyNameUpdated(event events.NameUpdated) {
}

func (a *Account) ApplyMoneyDeposited(event events.MoneyDeposited) {
	// No state change in aggregate for now, but could update balance if we had it
}

func (a *Account) ApplyMoneyWithdrawn(event events.MoneyWithdrawn) {
	// No state change in aggregate for now
}
