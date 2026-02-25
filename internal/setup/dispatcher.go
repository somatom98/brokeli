package setup

import (
	"github.com/somatom98/brokeli/internal/domain/account"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/pkg/event_store"
)

func TransactionDispatcher(es event_store.Store[*transaction.Transaction]) *transaction.Dispatcher {
	return transaction.NewDispatcher(es)
}

func AccountDispatcher(es event_store.Store[*account.Account]) *account.Dispatcher {
	return account.NewDispatcher(es)
}
