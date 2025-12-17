package create_transactions

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type Dispatcher interface {
	CreateExpense(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error
	CreateIncome(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error
	CreateTransfer(ctx context.Context, id uuid.UUID, fromAccountID uuid.UUID, fromCurrency values.Currency, fromAmount decimal.Decimal, toAccountID uuid.UUID, toCurrency values.Currency, toAmount decimal.Decimal, category, description string) error
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
