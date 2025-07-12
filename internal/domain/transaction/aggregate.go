package transaction

import (
	"github.com/gofrs/uuid"
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
		case events.Type_ExpenseCreated:
			t.ApplyExpenseCreated(record.Content().(events.ExpenseCreated))
		case events.Type_IncomeCreated:
			t.ApplyIncomeCreated(record.Content().(events.IncomeCreated))
		case events.Type_TransferCreated:
			t.ApplyTransferCreated(record.Content().(events.TransferCreated))
		}
	}

	return nil
}
