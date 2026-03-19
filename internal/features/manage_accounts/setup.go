package manage_accounts

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
	"github.com/somatom98/brokeli/internal/domain/projections/balance_updates"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type AccountDispatcher interface {
	Open(ctx context.Context, id uuid.UUID, name string, currency values.Currency, happenedAt time.Time) error
	UpdateName(ctx context.Context, id uuid.UUID, name string, happenedAt time.Time) error
	Deposit(ctx context.Context, id uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description, user string, happenedAt time.Time) error
	Withdraw(ctx context.Context, id uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description, user string, happenedAt time.Time) error
}

type Feature struct {
	httpHandler       *http.ServeMux
	accountsView      *accounts.Projection
	balanceUpdatesView *balance_updates.Projection
	accountDispatcher AccountDispatcher
}

func New(
	httpHandler *http.ServeMux,
	accountsView *accounts.Projection,
	balanceUpdatesView *balance_updates.Projection,
	accountDispatcher AccountDispatcher,
	transactionES event_store.Store[*transaction.Transaction],
) *Feature {
	f := &Feature{
		httpHandler:       httpHandler,
		accountsView:      accountsView,
		balanceUpdatesView: balanceUpdatesView,
		accountDispatcher: accountDispatcher,
	}

	transactionES.Subscribe(context.Background(), f.HandleRecord)

	return f
}

func (f *Feature) Setup(ctx context.Context) {
	f.httpHandler.HandleFunc("GET /api/accounts", f.handleGetAccounts)
	f.httpHandler.HandleFunc("POST /api/accounts", f.handleOpenAccount)
	f.httpHandler.HandleFunc("GET /api/accounts/{id}/balances", f.handleGetAccountBalances)
	f.httpHandler.HandleFunc("GET /api/accounts/{id}/distributions", f.handleGetAccountDistributions)
	f.httpHandler.HandleFunc("GET /api/balances", f.handleGetAllBalances)
	f.httpHandler.HandleFunc("POST /api/accounts/{id}/deposits", f.handleDeposit)
	f.httpHandler.HandleFunc("POST /api/accounts/{id}/withdrawals", f.handleWithdrawal)
}
