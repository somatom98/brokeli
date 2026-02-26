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

	processRecord := func(record event_store.Record) error {
		var err error
		switch record.Type() {
		case transaction_events.TypeMoneySpent:
			err = v.ApplyExpenseCreated(ctx, record.Content().(transaction_events.MoneySpent))
		case transaction_events.TypeMoneyReceived:
			err = v.ApplyIncomeCreated(ctx, record.Content().(transaction_events.MoneyReceived))
		case transaction_events.TypeReimbursementReceived:
			err = v.ApplyReimbursementReceived(ctx, record.Content().(transaction_events.ReimbursementReceived))
		case account_events.TypeOpened:
			err = v.ApplyAccountOpened(ctx, record.Content().(account_events.Opened))
		case account_events.TypeMoneyDeposited:
			err = v.ApplyMoneyDeposited(ctx, record.Content().(account_events.MoneyDeposited))
		case account_events.TypeMoneyWithdrawn:
			err = v.ApplyMoneyWithdrawn(ctx, record.Content().(account_events.MoneyWithdrawn))
		}
		return err
	}

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

				if err := processRecord(record); err != nil {
					errCh <- fmt.Errorf("apply transaction failed: %w", err)
					return
				}
			case record, ok := <-v.accountsCh:
				if !ok {
					errCh <- nil
					return
				}

				if err := processRecord(record); err != nil {
					errCh <- fmt.Errorf("apply account failed: %w", err)
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
