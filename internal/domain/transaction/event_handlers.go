package transaction

import (
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func (t *Transaction) ApplyExpenseCreated(e events.ExpenseCreated) {
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
