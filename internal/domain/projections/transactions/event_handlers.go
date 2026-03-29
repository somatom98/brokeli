package transactions

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/values"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

func (v *Projection) ApplyMoneySpent(ctx context.Context, idStr string, e transaction_events.MoneySpent) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Expense),
		Amount:          e.Amount.Neg(),
		Currency:        e.Currency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyReceived(ctx context.Context, idStr string, e transaction_events.MoneyReceived) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Income),
		Amount:          e.Amount,
		Currency:        e.Currency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyTransfered(ctx context.Context, idStr string, e transaction_events.MoneyTransfered) error {
	idSource := uuid.NewMD5(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s_source", idStr)))
	err := v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              idSource,
		AccountID:       e.FromAccountID,
		TransactionType: string(values.TransactionType_Transfer),
		Amount:          e.FromAmount.Neg(),
		Currency:        e.FromCurrency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
	if err != nil {
		return err
	}

	idDestination := uuid.NewMD5(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s_destination", idStr)))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              idDestination,
		AccountID:       e.ToAccountID,
		TransactionType: string(values.TransactionType_Transfer),
		Amount:          e.ToAmount,
		Currency:        e.ToCurrency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyReimbursementReceived(ctx context.Context, idStr string, e transaction_events.ReimbursementReceived) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Reimbursement),
		Amount:          e.Amount,
		Currency:        e.Currency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyDeposited(ctx context.Context, idStr string, e account_events.MoneyDeposited) error {
	if e.User == "system" {
		return nil
	}

	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Deposit),
		Amount:          e.Amount,
		Currency:        e.Currency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyWithdrawn(ctx context.Context, idStr string, e account_events.MoneyWithdrawn) error {
	if e.User == "system" {
		return nil
	}

	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Withdrawal),
		Amount:          e.Amount.Neg(),
		Currency:        e.Currency,
		Category:        e.Category,
		Description:     e.Description,
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyInvested(ctx context.Context, idStr string, e transaction_events.MoneyInvested) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))

	if e.PriceCurrency == e.FeeCurrency {
		amount := e.Units.Mul(e.Price).Add(e.Fee)
		return v.repository.CreateTransaction(ctx, TransactionRecord{
			ID:              id,
			AccountID:       e.AccountID,
			TransactionType: string(values.TransactionType_Investment),
			Amount:          amount.Neg(), // Money is leaving liquidity
			Currency:        e.PriceCurrency,
			Category:        "Investments",
			Description:     e.Ticker,
			HappenedAt:      e.HappenedAt,
		})
	}

	// Currencies are different, create two records
	priceAmount := e.Units.Mul(e.Price)
	err := v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Investment),
		Amount:          priceAmount.Neg(),
		Currency:        e.PriceCurrency,
		Category:        "Investments",
		Description:     e.Ticker,
		HappenedAt:      e.HappenedAt,
	})
	if err != nil {
		return err
	}

	idFee := uuid.NewMD5(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s_fee", idStr)))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              idFee,
		AccountID:       e.AccountID,
		TransactionType: string(values.TransactionType_Investment),
		Amount:          e.Fee.Neg(),
		Currency:        e.FeeCurrency,
		Category:        "Investments",
		Description:     fmt.Sprintf("%s (Fee)", e.Ticker),
		HappenedAt:      e.HappenedAt,
	})
}

