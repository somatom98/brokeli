package setup

import (
	"context"

	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/pkg/event_store"
)

func AccountsProjection(
	ctx context.Context,
	transactionES event_store.Store[*transaction.Transaction],
) *accounts.Projection {
	accountsProjection := accounts.New(transactionES)
	_ = accountsProjection.Update(ctx)
	return accountsProjection
}
