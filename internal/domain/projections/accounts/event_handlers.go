package accounts

import (
	"context"

	"github.com/google/uuid"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

func (v *Projection) ApplyAccountCreated(ctx context.Context, id uuid.UUID, e account_events.Created) error {
	return v.repository.CreateAccount(ctx, id, e.Time)
}

func (v *Projection) ApplyAccountClosed(ctx context.Context, id uuid.UUID, e account_events.AccountClosed) error {
	return v.repository.CloseAccount(ctx, id, e.Time)
}

func (v *Projection) ApplyAccountDeposited(ctx context.Context, id uuid.UUID, e account_events.MoneyDeposited) error {
	return v.repository.UpdateAccountBalance(ctx, id, e.Amount, e.Currency)
}

func (v *Projection) ApplyExpenseCreated(ctx context.Context, e transaction_events.MoneySpent) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount.Neg(), e.Currency)
}

func (v *Projection) ApplyIncomeCreated(ctx context.Context, e transaction_events.MoneyReceived) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyExpectedReimbursementSet(ctx context.Context, e transaction_events.ExpectedReimbursementSet) error {
	return v.repository.SetExpectedReimbursement(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyReimbursementReceived(ctx context.Context, e transaction_events.ReimbursementReceived) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
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
