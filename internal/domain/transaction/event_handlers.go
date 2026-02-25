package transaction

import (
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func (t *Transaction) ApplyExpectedReimbursementSet(e events.ExpectedReimbursementSet) {
	t.State = State_Created
	t.Type = values.TransactionType_ExpectedReimbursement
	t.Description = ""
}

func (t *Transaction) ApplyExpenseCreated(e events.MoneySpent) {
	t.State = State_Created
	t.Type = values.TransactionType_Expense
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

func (t *Transaction) ApplyIncomeCreated(e events.MoneyReceived) {
	t.State = State_Created
	t.Type = values.TransactionType_Income
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

func (t *Transaction) ApplyTransferCreated(e events.MoneyTransfered) {
	t.State = State_Created
	t.Type = values.TransactionType_Transfer
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

func (t *Transaction) ApplyReimbursementReceived(e events.ReimbursementReceived) {
	t.State = State_Created
	t.Type = values.TransactionType_Reimbursement
	entry := values.Entry{
		AccountID: e.AccountID,
		Currency:  e.Currency,
		Amount:    e.Amount,
		Side:      values.Side_Credit,
	}
	t.Entries = append(t.Entries, entry)
	t.Description = e.From
}

func (t *Transaction) ApplyMoneyDeposited(e events.MoneyDeposited) {
	t.State = State_Created
	t.Type = values.TransactionType_Deposit
	entry := values.Entry{
		AccountID: e.AccountID,
		Currency:  e.Currency,
		Amount:    e.Amount,
		Side:      values.Side_Credit,
	}
	t.Entries = append(t.Entries, entry)
}

func (t *Transaction) ApplyMoneyWithdrawn(e events.MoneyWithdrawn) {
	t.State = State_Created
	t.Type = values.TransactionType_Withdrawal
	entry := values.Entry{
		AccountID: e.AccountID,
		Currency:  e.Currency,
		Amount:    e.Amount,
		Side:      values.Side_Debit,
	}
	t.Entries = append(t.Entries, entry)
}
