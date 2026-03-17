package transactions

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

const (
	TypeDebit  = "DEBIT"
	TypeCredit = "CREDIT"
)

func (v *Projection) ApplyMoneySpent(ctx context.Context, idStr string, e transaction_events.MoneySpent) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: TypeDebit,
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
		TransactionType: TypeCredit,
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
		TransactionType: TypeDebit,
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
		TransactionType: TypeCredit,
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
		TransactionType: TypeCredit,
		Amount:          e.Amount,
		Currency:        e.Currency,
		Category:        "Reimbursement",
		Description:     fmt.Sprintf("Reimbursement from %s", e.From),
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyDeposited(ctx context.Context, idStr string, e account_events.MoneyDeposited) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: TypeCredit,
		Amount:          e.Amount,
		Currency:        e.Currency,
		Category:        "Deposit",
		Description:     fmt.Sprintf("Deposit from %s", e.User),
		HappenedAt:      e.HappenedAt,
	})
}

func (v *Projection) ApplyMoneyWithdrawn(ctx context.Context, idStr string, e account_events.MoneyWithdrawn) error {
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))
	return v.repository.CreateTransaction(ctx, TransactionRecord{
		ID:              id,
		AccountID:       e.AccountID,
		TransactionType: TypeDebit,
		Amount:          e.Amount.Neg(),
		Currency:        e.Currency,
		Category:        "Withdrawal",
		Description:     fmt.Sprintf("Withdrawal by %s", e.User),
		HappenedAt:      e.HappenedAt,
	})
}
