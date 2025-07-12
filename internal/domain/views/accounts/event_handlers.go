package accounts

import (
	"context"

	"github.com/somatom98/brokeli/internal/domain/transaction/events"
)

func (v *View) ApplyExpenseCreated(ctx context.Context, e events.ExpenseCreated) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount.Neg(), e.Currency)
}

func (v *View) ApplyIncomeCreated(ctx context.Context, e events.IncomeCreated) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *View) ApplyTransferCreated(ctx context.Context, e events.TransferCreated) error {
	err := v.repository.UpdateAccountBalance(ctx, e.FromAccountID, e.FromAmount.Neg(), e.FromCurrency)
	if err != nil {
		return err
	}

	err = v.repository.UpdateAccountBalance(ctx, e.ToAccountID, e.ToAmount, e.ToCurrency)
	if err != nil {
		return err
	}

	return nil
}
