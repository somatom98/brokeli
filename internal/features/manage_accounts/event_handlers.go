package manage_accounts

import (
	"context"
	"fmt"

	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/pkg/event_store"
)

func (f *Feature) HandleRecord(ctx context.Context, record event_store.Record) error {
	switch record.Type() {
	case transaction_events.TypeMoneyTransfered:
		event := record.Content().(transaction_events.MoneyTransfered)

		return f.handleMoneyTransfered(ctx, event)
	}
	return nil
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
