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
	t.Category = e.Category
	t.Description = e.Description
}

func (t *Transaction) ApplyInvestmentCreated(e events.MoneyInvested) {
	t.State = State_Created
	t.Type = values.TransactionType_Investment

	if e.PriceCurrency == e.FeeCurrency {
		amount := e.Units.Mul(e.Price).Add(e.Fee)
		entry := values.Entry{
			AccountID: e.AccountID,
			Currency:  e.PriceCurrency,
			Amount:    amount,
			Side:      values.Side_Debit,
		}
		t.Entries = append(t.Entries, entry)
	} else {
		priceAmount := e.Units.Mul(e.Price)
		priceEntry := values.Entry{
			AccountID: e.AccountID,
			Currency:  e.PriceCurrency,
			Amount:    priceAmount,
			Side:      values.Side_Debit,
		}
		feeEntry := values.Entry{
			AccountID: e.AccountID,
			Currency:  e.FeeCurrency,
			Amount:    e.Fee,
			Side:      values.Side_Debit,
		}
		t.Entries = append(t.Entries, priceEntry, feeEntry)
	}

	t.Category = "Investments"
	t.Description = e.Ticker
}
