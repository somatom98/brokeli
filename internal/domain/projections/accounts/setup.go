package accounts

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

type Repository interface {
	CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error
	CloseAccount(ctx context.Context, id uuid.UUID, closedAt time.Time) error
	UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error
	GetAll(ctx context.Context) (map[uuid.UUID]Account, error)
}

type Projection struct {
	repository     Repository
	transactionsCh <-chan event_store.Record
	accountsCh     <-chan event_store.Record
}

func New(
	transactionES event_store.Store[*transaction.Transaction],
	accountES event_store.Store[*account.Account],
	repository Repository,
) *Projection {
	return &Projection{
		repository:     repository,
		transactionsCh: transactionES.Subscribe(context.Background()),
		accountsCh:     accountES.Subscribe(context.Background()),
	}
}

func (v *Projection) Update(ctx context.Context) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		transactionsClosed := false
		accountsClosed := false

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case record, ok := <-v.transactionsCh:
				if !ok {
					transactionsClosed = true
					if accountsClosed {
						errCh <- nil
						return
					}
					continue
				}

				var err error
				switch record.Type() {
				case transaction_events.Type_MoneySpent:
					err = v.ApplyExpenseCreated(ctx, record.Content().(transaction_events.MoneySpent))
				case transaction_events.Type_MoneyReceived:
					err = v.ApplyIncomeCreated(ctx, record.Content().(transaction_events.MoneyReceived))
				case transaction_events.Type_ReimbursementReceived:
					err = v.ApplyReimbursementReceived(ctx, record.Content().(transaction_events.ReimbursementReceived))
				case transaction_events.Type_MoneyTransfered:
					err = v.ApplyTransferCreated(ctx, record.Content().(transaction_events.MoneyTransfered))
				}

				if err != nil {
					errCh <- fmt.Errorf("apply failed: %w", err)
					return
				}
			case record, ok := <-v.accountsCh:
				if !ok {
					accountsClosed = true
					if transactionsClosed {
						errCh <- nil
						return
					}
					continue
				}

				var err error
				switch record.Type() {
				case account_events.Type_Created:
					err = v.ApplyAccountCreated(ctx, record.AggregateID, record.Content().(account_events.Created))
				case account_events.Type_AccountClosed:
					err = v.ApplyAccountClosed(ctx, record.AggregateID, record.Content().(account_events.AccountClosed))
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
