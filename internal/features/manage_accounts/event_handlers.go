package manage_accounts

import (
	"context"
	"fmt"

	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
)

func (f *Feature) Listen(ctx context.Context) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case record, ok := <-f.transactionsCh:
				if !ok {
					errCh <- nil
					return
				}

				switch record.Type() {
				case transaction_events.TypeMoneyTransfered:
					event := record.Content().(transaction_events.MoneyTransfered)

					err := f.handleMoneyTransfered(ctx, event)
					if err != nil {
						errCh <- err
						return
					}
				}
			}
		}
	}()

	return errCh
}

func (f *Feature) handleMoneyTransfered(ctx context.Context, event transaction_events.MoneyTransfered) error {
	err := f.accountDispatcher.Withdraw(
		ctx,
		event.FromAccountID,
		event.FromCurrency,
		event.FromAmount,
		"system",
	)
	if err != nil {
		return fmt.Errorf("failed to process transfer withdrawal: %w", err)
	}

	err = f.accountDispatcher.Deposit(
		ctx,
		event.ToAccountID,
		event.ToCurrency,
		event.ToAmount,
		"system",
	)
	if err != nil {
		return fmt.Errorf("failed to process transfer deposit: %w", err)
	}

	return nil
}
