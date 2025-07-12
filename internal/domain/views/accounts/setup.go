package accounts

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type Repository interface {
	UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error
	GetAll(ctx context.Context) ([]Account, error)
}

type View struct {
	repository     Repository
	transactionsCh <-chan event_store.Record
}

func New(transactionES event_store.Store[*transaction.Transaction]) *View {
	return &View{
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
				case transaction_events.Type_ExpenseCreated:
					err = v.ApplyExpenseCreated(ctx, record.Content().(transaction_events.ExpenseCreated))
				case transaction_events.Type_IncomeCreated:
					err = v.ApplyIncomeCreated(ctx, record.Content().(transaction_events.IncomeCreated))
				case transaction_events.Type_TransferCreated:
					err = v.ApplyTransferCreated(ctx, record.Content().(transaction_events.TransferCreated))
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

func (v *View) GetAll(ctx context.Context) ([]Account, error) {
	return v.repository.GetAll(ctx)
}
