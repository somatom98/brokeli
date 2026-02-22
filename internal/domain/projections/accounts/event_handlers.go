package accounts

import (
	"context"

	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

func (v *Projection) ApplyExpenseCreated(ctx context.Context, e transaction_events.MoneySpent) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount.Neg(), e.Currency)
}

func (v *Projection) ApplyIncomeCreated(ctx context.Context, e transaction_events.MoneyReceived) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyReimbursementReceived(ctx context.Context, e transaction_events.ReimbursementReceived) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyMoneyDeposited(ctx context.Context, e transaction_events.MoneyDeposited) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyMoneyWithdrawn(ctx context.Context, e transaction_events.MoneyWithdrawn) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount.Neg(), e.Currency)
}

func (v *Projection) ApplyTransferCreated(ctx context.Context, e transaction_events.MoneyTransfered) error {
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
