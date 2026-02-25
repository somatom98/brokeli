package accounts

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	repository     Repository
	transactionsCh <-chan event_store.Record
}

func New(
	transactionES event_store.Store[*transaction.Transaction],
	repository Repository,
) *Projection {
	return &Projection{
		repository:     repository,
		transactionsCh: transactionES.Subscribe(context.Background()),
	}
}

func (v *Projection) Update(ctx context.Context) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case record, ok := <-v.transactionsCh:
				if !ok {
					errCh <- nil
					return
				}

				var err error
				switch record.Type() {
				case transaction_events.TypeMoneySpent:
					err = v.ApplyExpenseCreated(ctx, record.Content().(transaction_events.MoneySpent))
				case transaction_events.TypeMoneyReceived:
					err = v.ApplyIncomeCreated(ctx, record.Content().(transaction_events.MoneyReceived))
				case transaction_events.TypeReimbursementReceived:
					err = v.ApplyReimbursementReceived(ctx, record.Content().(transaction_events.ReimbursementReceived))
				case transaction_events.TypeMoneyTransfered:
					err = v.ApplyTransferCreated(ctx, record.Content().(transaction_events.MoneyTransfered))
				case transaction_events.TypeMoneyDeposited:
					err = v.ApplyMoneyDeposited(ctx, record.Content().(transaction_events.MoneyDeposited))
				case transaction_events.TypeMoneyWithdrawn:
					err = v.ApplyMoneyWithdrawn(ctx, record.Content().(transaction_events.MoneyWithdrawn))
				}

				if err != nil {
					errCh <- fmt.Errorf("apply failed: %w", err)
					return
				}
			}
		}
	}()

	return errCh
}

func (v *Projection) GetAll(ctx context.Context) (map[uuid.UUID]Account, error) {
	return v.repository.GetAll(ctx)
}
