package transaction

import (
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

func (t *Transaction) Hydrate(records []event_store.Record) error {
	for _, record := range records {
		switch record.Event.Type() {
		case "ExpenseCreated":
			t.ApplyExpenseCreated(record.Event.Content().(events.ExpenseCreated))
		case "IncomeCreated":
			t.ApplyIncomeCreated(record.Event.Content().(events.IncomeCreated))
		case "TransferCreated":
			t.ApplyTransferCreated(record.Event.Content().(events.TransferCreated))
		}
	}

	return nil
}

func (t *Transaction) ApplyExpenseCreated(e events.ExpenseCreated) {
	t.State = State_Created
	entry := values.Entry{
		AccountID: e.AccountID,
		Currency:  e.Currency,
		Amount:    e.Amount,
		Side:      values.Side_Debit,
	}
	t.Entries = append(t.Entries, entry)
	t.Category = e.Category
	t.Description = e.Description
}

func (t *Transaction) ApplyIncomeCreated(e events.IncomeCreated) {
	t.State = State_Created
	entry := values.Entry{
		AccountID: e.AccountID,
		Currency:  e.Currency,
		Amount:    e.Amount,
		Side:      values.Side_Credit,
	}
	t.Entries = append(t.Entries, entry)
	t.Category = e.Category
	t.Description = e.Description
}

func (t *Transaction) ApplyTransferCreated(e events.TransferCreated) {
	t.State = State_Created
	from := values.Entry{
		AccountID: e.FromAccountID,
		Currency:  e.FromCurrency,
		Amount:    e.FromAmount,
		Side:      values.Side_Debit,
	}
	to := values.Entry{
		AccountID: e.ToAccountID,
		Currency:  e.ToCurrency,
		Amount:    e.ToAmount,
		Side:      values.Side_Credit,
	}
	t.Entries = append(t.Entries, from, to)
	t.Category = e.Category
	t.Description = e.Description
}
