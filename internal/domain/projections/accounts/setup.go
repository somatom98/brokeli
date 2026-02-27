package accounts

import (
	"context"
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

type Repository interface {
	CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error
	CloseAccount(ctx context.Context, id uuid.UUID, closedAt time.Time) error
	UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error
	GetAll(ctx context.Context) (map[uuid.UUID]Account, error)
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
	switch record.Type() {
	case transaction_events.TypeMoneySpent:
		return v.ApplyExpenseCreated(ctx, record.Content().(transaction_events.MoneySpent))
	case transaction_events.TypeMoneyReceived:
		return v.ApplyIncomeCreated(ctx, record.Content().(transaction_events.MoneyReceived))
	case transaction_events.TypeReimbursementReceived:
		return v.ApplyReimbursementReceived(ctx, record.Content().(transaction_events.ReimbursementReceived))
	case account_events.TypeOpened:
		return v.ApplyAccountOpened(ctx, record.Content().(account_events.Opened))
	case account_events.TypeMoneyDeposited:
		return v.ApplyMoneyDeposited(ctx, record.Content().(account_events.MoneyDeposited))
	case account_events.TypeMoneyWithdrawn:
		return v.ApplyMoneyWithdrawn(ctx, record.Content().(account_events.MoneyWithdrawn))
	}
	return nil
}

func (v *Projection) GetAll(ctx context.Context) (map[uuid.UUID]Account, error) {
	return v.repository.GetAll(ctx)
}
