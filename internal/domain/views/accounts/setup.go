package accounts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type Repository interface {
	UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error
	GetAll(ctx context.Context) (map[uuid.UUID]Account, error)
}

type View struct {
	repository     Repository
	transactionsCh <-chan event_store.Record
}

func New(
	transactionES event_store.Store[*transaction.Transaction],
) *View {
	return &View{
		repository:     NewInMemoryRepository(),
		transactionsCh: transactionES.Subscribe(context.Background()),
	}
}

func (v *View) Update(ctx context.Context) <-chan error {
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
				case transaction_events.Type_MoneySpent:
					err = v.ApplyExpenseCreated(ctx, record.Content().(transaction_events.MoneySpent))
				case transaction_events.Type_MoneyReceived:
					err = v.ApplyIncomeCreated(ctx, record.Content().(transaction_events.MoneyReceived))
				case transaction_events.Type_MoneyTransfered:
					err = v.ApplyTransferCreated(ctx, record.Content().(transaction_events.MoneyTransfered))
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

func (v *View) GetAll(ctx context.Context) (map[uuid.UUID]Account, error) {
	return v.repository.GetAll(ctx)
}
