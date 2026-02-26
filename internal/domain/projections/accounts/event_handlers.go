package accounts

import (
	"context"
	"time"

	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

func (v *Projection) ApplyAccountOpened(ctx context.Context, e account_events.Opened) error {
	return v.repository.CreateAccount(ctx, e.AccountID, time.Now())
}

func (v *Projection) ApplyExpenseCreated(ctx context.Context, e transaction_events.MoneySpent) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount.Neg(), e.Currency)
}

func (v *Projection) ApplyIncomeCreated(ctx context.Context, e transaction_events.MoneyReceived) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyReimbursementReceived(ctx context.Context, e transaction_events.ReimbursementReceived) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyMoneyDeposited(ctx context.Context, e account_events.MoneyDeposited) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount, e.Currency)
}

func (v *Projection) ApplyMoneyWithdrawn(ctx context.Context, e account_events.MoneyWithdrawn) error {
	return v.repository.UpdateAccountBalance(ctx, e.AccountID, e.Amount.Neg(), e.Currency)
}
