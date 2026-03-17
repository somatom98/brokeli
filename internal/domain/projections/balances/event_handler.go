package balances

import (
	"context"

	"github.com/google/uuid"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

const userSystem = "system"

const (
	originTransaction = "transaction"
	originMovement    = "movement"
)

func (v *Projection) ApplyExpenseCreated(ctx context.Context, id uuid.UUID, e transaction_events.MoneySpent) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount.Neg(), userSystem, e.HappenedAt, originTransaction)
}

func (v *Projection) ApplyIncomeCreated(ctx context.Context, id uuid.UUID, e transaction_events.MoneyReceived) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount, userSystem, e.HappenedAt, originTransaction)
}

func (v *Projection) ApplyReimbursementReceived(ctx context.Context, id uuid.UUID, e transaction_events.ReimbursementReceived) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount, userSystem, e.HappenedAt, originTransaction)
}

func (v *Projection) ApplyMoneyDeposited(ctx context.Context, id uuid.UUID, e account_events.MoneyDeposited) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount, e.User, e.HappenedAt, originMovement)
}

func (v *Projection) ApplyMoneyWithdrawn(ctx context.Context, id uuid.UUID, e account_events.MoneyWithdrawn) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount.Neg(), e.User, e.HappenedAt, originMovement)
}
