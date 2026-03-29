package balance_updates

import (
	"context"
	"fmt"

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
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount.Neg(), userSystem, e.HappenedAt, originTransaction, BalanceTypeLiquidity)
}

func (v *Projection) ApplyIncomeCreated(ctx context.Context, id uuid.UUID, e transaction_events.MoneyReceived) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount, userSystem, e.HappenedAt, originTransaction, BalanceTypeLiquidity)
}

func (v *Projection) ApplyReimbursementReceived(ctx context.Context, id uuid.UUID, e transaction_events.ReimbursementReceived) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount, userSystem, e.HappenedAt, originTransaction, BalanceTypeLiquidity)
}

func (v *Projection) ApplyMoneyDeposited(ctx context.Context, id uuid.UUID, e account_events.MoneyDeposited) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount, e.User, e.HappenedAt, originMovement, BalanceTypeLiquidity)
}

func (v *Projection) ApplyMoneyWithdrawn(ctx context.Context, id uuid.UUID, e account_events.MoneyWithdrawn) error {
	return v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.Currency, e.Amount.Neg(), e.User, e.HappenedAt, originMovement, BalanceTypeLiquidity)
}

func (v *Projection) ApplyInvestmentCreated(ctx context.Context, id uuid.UUID, e transaction_events.MoneyInvested) error {
	priceAmount := e.Units.Mul(e.Price)

	// 1. Withdrawal from liquidity (Asset cost)
	err := v.repository.InsertBalanceUpdate(ctx, id, e.AccountID, e.PriceCurrency, priceAmount.Neg(), userSystem, e.HappenedAt, originTransaction, BalanceTypeLiquidity)
	if err != nil {
		return err
	}

	// 2. Withdrawal from liquidity (Fee)
	idFee := uuid.NewMD5(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s_fee", id.String())))
	err = v.repository.InsertBalanceUpdate(ctx, idFee, e.AccountID, e.FeeCurrency, e.Fee.Neg(), userSystem, e.HappenedAt, originTransaction, BalanceTypeLiquidity)
	if err != nil {
		return err
	}

	// 3. Deposit into investment balance (Asset cost)
	idInvestment := uuid.NewMD5(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s_investment", id.String())))
	return v.repository.InsertBalanceUpdate(ctx, idInvestment, e.AccountID, e.PriceCurrency, priceAmount, userSystem, e.HappenedAt, originTransaction, BalanceTypeInvestment)
}
