package transaction

import (
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
			t.ApplyExpectedReimbursementSet(record.Content().(events.ExpectedReimbursementSet))
		case events.Type_MoneySpent:
			t.ApplyExpenseCreated(record.Content().(events.MoneySpent))
		case events.Type_MoneyReceived:
			t.ApplyIncomeCreated(record.Content().(events.MoneyReceived))
		case events.Type_MoneyTransfered:
			t.ApplyTransferCreated(record.Content().(events.MoneyTransfered))
		case events.Type_ReimbursementReceived:
			t.ApplyReimbursementReceived(record.Content().(events.ReimbursementReceived))
		}
	}

	return nil
}
