package create_transactions

import (
	"context"
	"net/http"

	"github.com/somatom98/brokeli/internal/domain/transaction"
)

type Dispatcher interface {
	CreateExpense(ctx context.Context, cmd transaction.CreateExpense) error
	CreateIncome(ctx context.Context, cmd transaction.CreateIncome) error
	CreateTransfer(ctx context.Context, cmd transaction.CreateTransfer) error
}

type Feature struct {
	httpHandler *http.ServeMux
	dispatcher  Dispatcher
}

func New(
	httpHandler *http.ServeMux,
	api Dispatcher,
) *Feature {
	return &Feature{
		httpHandler: httpHandler,
		dispatcher:  api,
	}
}

func (f *Feature) Setup() {
	f.httpHandler.HandleFunc("POST /api/expenses", f.handleCreateExpense)
	f.httpHandler.HandleFunc("POST /api/incomes", f.handleCreateIncome)
	f.httpHandler.HandleFunc("POST /api/transfers", f.handleCreateTransfer)
}
