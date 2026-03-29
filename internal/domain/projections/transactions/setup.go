package transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type TransactionRecord struct {
	ID              uuid.UUID       `json:"id"`
	AccountID       uuid.UUID       `json:"account_id"`
	TransactionType string          `json:"transaction_type"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        values.Currency `json:"currency"`
	Category        string          `json:"category"`
	Description     string          `json:"description"`
	HappenedAt      time.Time       `json:"happened_at"`
	SystemTotalRate decimal.Decimal `json:"system_total_rate"`
}

type ListTransactionsParams struct {
	StartDate       *time.Time
	EndDate         *time.Time
	AccountIDs      []uuid.UUID
	TransactionType *string
}

type ListTransactionsPaginatedParams struct {
	ListTransactionsParams
	Limit  int32
	Offset int32
}

type PaginatedTransactions struct {
	Transactions []TransactionRecord `json:"transactions"`
	TotalCount   int64               `json:"total_count"`
}

type Repository interface {
	CreateTransaction(ctx context.Context, tx TransactionRecord) error
	ListTransactions(ctx context.Context, params ListTransactionsParams) ([]TransactionRecord, error)
	ListTransactionsPaginated(ctx context.Context, params ListTransactionsPaginatedParams) (PaginatedTransactions, error)
	ListCategories(ctx context.Context) ([]string, error)
}

type Projection struct {
	repository Repository
}

func New(
	transactionES event_store.Store[*transaction.Transaction],
	accountES event_store.Store[*account.Account],
	repository Repository,
) *Projection {
	p := &Projection{
		repository: repository,
	}

	transactionES.Subscribe(context.Background(), p.HandleRecord)
	accountES.Subscribe(context.Background(), p.HandleRecord)

	return p
}

func (v *Projection) HandleRecord(ctx context.Context, record event_store.Record) error {
	var aggregateType string
	switch record.Type() {
	case transaction_events.TypeMoneySpent, transaction_events.TypeMoneyReceived, transaction_events.TypeMoneyTransfered, transaction_events.TypeReimbursementReceived, transaction_events.TypeMoneyInvested:
		aggregateType = "Transaction"
	case account_events.TypeMoneyDeposited, account_events.TypeMoneyWithdrawn:
		aggregateType = "Account"
	default:
		return nil
	}

	idStr := fmt.Sprintf("%s_%s_%d", aggregateType, record.AggregateID.String(), record.Version)

	switch record.Type() {
	case transaction_events.TypeMoneySpent:
		return v.ApplyMoneySpent(ctx, idStr, record.Content().(transaction_events.MoneySpent))
	case transaction_events.TypeMoneyReceived:
		return v.ApplyMoneyReceived(ctx, idStr, record.Content().(transaction_events.MoneyReceived))
	case transaction_events.TypeMoneyTransfered:
		return v.ApplyMoneyTransfered(ctx, idStr, record.Content().(transaction_events.MoneyTransfered))
	case transaction_events.TypeReimbursementReceived:
		return v.ApplyReimbursementReceived(ctx, idStr, record.Content().(transaction_events.ReimbursementReceived))
	case transaction_events.TypeMoneyInvested:
		return v.ApplyMoneyInvested(ctx, idStr, record.Content().(transaction_events.MoneyInvested))
	case account_events.TypeMoneyDeposited:
		return v.ApplyMoneyDeposited(ctx, idStr, record.Content().(account_events.MoneyDeposited))
	case account_events.TypeMoneyWithdrawn:
		return v.ApplyMoneyWithdrawn(ctx, idStr, record.Content().(account_events.MoneyWithdrawn))
	}
	return nil
}

func (v *Projection) ListTransactions(ctx context.Context, params ListTransactionsParams) ([]TransactionRecord, error) {
	return v.repository.ListTransactions(ctx, params)
}

func (v *Projection) ListTransactionsPaginated(ctx context.Context, params ListTransactionsPaginatedParams) (PaginatedTransactions, error) {
	return v.repository.ListTransactionsPaginated(ctx, params)
}

func (v *Projection) ListCategories(ctx context.Context) ([]string, error) {
	return v.repository.ListCategories(ctx)
}
