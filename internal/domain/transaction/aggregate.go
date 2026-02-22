package transaction

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

const (
	Version = 0
)

type State int

const (
	State_Created State = iota
	State_Deleted
)

type Transaction struct {
	ID          uuid.UUID
	State       State
	Type        values.TransactionType
	Entries     []values.Entry
	Category    string
	Description string
}

func New(id uuid.UUID) *Transaction {
	return &Transaction{
		ID:      id,
		Entries: []values.Entry{},
	}
}

func (t *Transaction) Hydrate(records []event_store.Record) error {
	for _, record := range records {
		switch record.Type() {
		case events.Type_ExpectedReimbursementSet:
			event, err := event_store.DecodeEvent[events.ExpectedReimbursementSet](record.Content())
			if err != nil {
				return fmt.Errorf("decode ExpectedReimbursementSet event: %w", err)
			}
			t.ApplyExpectedReimbursementSet(event)
		case events.Type_MoneySpent:
			event, err := event_store.DecodeEvent[events.MoneySpent](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneySpent event: %w", err)
			}
			t.ApplyExpenseCreated(event)
		case events.Type_MoneyReceived:
			event, err := event_store.DecodeEvent[events.MoneyReceived](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyReceived event: %w", err)
			}
			t.ApplyIncomeCreated(event)
		case events.Type_MoneyTransfered:
			event, err := event_store.DecodeEvent[events.MoneyTransfered](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyTransfered event: %w", err)
			}
			t.ApplyTransferCreated(event)
		case events.Type_MoneyDeposited:
			event, err := event_store.DecodeEvent[events.MoneyDeposited](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyDeposited event: %w", err)
			}
			t.ApplyMoneyDeposited(event)
		case events.Type_MoneyWithdrawn:
			event, err := event_store.DecodeEvent[events.MoneyWithdrawn](record.Content())
			if err != nil {
				return fmt.Errorf("decode MoneyWithdrawn event: %w", err)
			}
			t.ApplyMoneyWithdrawn(event)
		case events.Type_ReimbursementReceived:
			event, err := event_store.DecodeEvent[events.ReimbursementReceived](record.Content())
			if err != nil {
				return fmt.Errorf("decode ReimbursementReceived event: %w", err)
			}
			t.ApplyReimbursementReceived(event)
		}
	}

	return nil
}
