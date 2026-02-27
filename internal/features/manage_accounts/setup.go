package manage_accounts

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type AccountDispatcher interface {
	Open(ctx context.Context, id uuid.UUID, name string, currency values.Currency) error
	UpdateName(ctx context.Context, id uuid.UUID, name string) error
	Deposit(ctx context.Context, id uuid.UUID, currency values.Currency, amount decimal.Decimal, user string) error
	Withdraw(ctx context.Context, id uuid.UUID, currency values.Currency, amount decimal.Decimal, user string) error
}

type Feature struct {
	httpHandler       *http.ServeMux
	accountsView      *accounts.Projection
	accountDispatcher AccountDispatcher
}

func New(
	httpHandler *http.ServeMux,
	accountsView *accounts.Projection,
	accountDispatcher AccountDispatcher,
	transactionES event_store.Store[*transaction.Transaction],
) *Feature {
	f := &Feature{
		httpHandler:       httpHandler,
		accountsView:      accountsView,
		accountDispatcher: accountDispatcher,
	}

	transactionES.Subscribe(context.Background(), f.HandleRecord)

	return f
}

func (f *Feature) Setup(ctx context.Context) {
	f.httpHandler.HandleFunc("GET /api/accounts", f.handleGetAccounts)
}
