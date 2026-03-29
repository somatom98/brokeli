package manage_transactions

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/transactions"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type Dispatcher interface {
	RegisterExpense(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string, happenedAt time.Time) error
	RegisterIncome(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string, happenedAt time.Time) error
	RegisterTransfer(ctx context.Context, id uuid.UUID, fromAccountID uuid.UUID, fromCurrency values.Currency, fromAmount decimal.Decimal, toAccountID uuid.UUID, toCurrency values.Currency, toAmount decimal.Decimal, category, description string, happenedAt time.Time) error
	RegisterReimbursement(ctx context.Context, id uuid.UUID, accountID uuid.UUID, from string, currency values.Currency, amount decimal.Decimal, category string, description string, happenedAt time.Time) error
	RegisterInvestment(ctx context.Context, id uuid.UUID, accountID uuid.UUID, ticker string, units decimal.Decimal, price decimal.Decimal, priceCurrency values.Currency, fee decimal.Decimal, feeCurrency values.Currency, happenedAt time.Time) error
	SetExpectedReimbursement(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, happenedAt time.Time) error
}

type Feature struct {
	httpHandler      *http.ServeMux
	dispatcher       Dispatcher
	transactionsView *transactions.Projection
}

func New(
	httpHandler *http.ServeMux,
	api Dispatcher,
	transactionsView *transactions.Projection,
) *Feature {
	return &Feature{
		httpHandler:      httpHandler,
		dispatcher:       api,
		transactionsView: transactionsView,
	}
}

func (f *Feature) Setup() {
	f.httpHandler.HandleFunc("GET /api/transactions", f.handleGetTransactions)
	f.httpHandler.HandleFunc("POST /api/expenses", f.handleRegisterExpense)
	f.httpHandler.HandleFunc("POST /api/incomes", f.handleRegisterIncome)
	f.httpHandler.HandleFunc("POST /api/transfers", f.handleRegisterTransfer)
	f.httpHandler.HandleFunc("POST /api/investments", f.handleRegisterInvestment)
	f.httpHandler.HandleFunc("POST /api/{transaction_id}/reimbursement", f.handleRegisterReimbursement)
	f.httpHandler.HandleFunc("POST /api/{transaction_id}/expected-reimbursements", f.handleSetExpectedReimbursement)
}
