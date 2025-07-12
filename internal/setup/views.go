package setup

import (
	"context"

	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/views/accounts"
	"github.com/somatom98/brokeli/pkg/event_store"
)

func AccountsView(
	ctx context.Context,
	transactionES event_store.Store[*transaction.Transaction],
) *accounts.View {
	accountsView := accounts.New(transactionES)
	_ = accountsView.Update(ctx)
	return accountsView
}
