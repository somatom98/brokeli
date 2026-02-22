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
	repository accounts.Repository,
) *accounts.Projection {
	accountsProjection := accounts.New(transactionES, repository)
	_ = accountsProjection.Update(ctx)
	return accountsProjection
}
