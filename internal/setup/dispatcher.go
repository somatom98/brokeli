package setup

import "github.com/somatom98/brokeli/internal/domain/transaction"

func Dispatcher() *transaction.Dispatcher {
	return transaction.NewDispatcher()
}
