package manage_accounts

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
)

type Dispatcher interface {
	CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error
	CloseAccount(ctx context.Context, id uuid.UUID, at time.Time) error
}

type Feature struct {
	httpHandler  *http.ServeMux
	accountsView *accounts.Projection
	dispatcher   Dispatcher
}

func New(
	httpHandler *http.ServeMux,
	accountsView *accounts.Projection,
	dispatcher Dispatcher,
) *Feature {
	return &Feature{
		httpHandler:  httpHandler,
		accountsView: accountsView,
		dispatcher:   dispatcher,
	}
}

func (f *Feature) Setup() {
	f.httpHandler.HandleFunc("GET /api/accounts", f.handleGetAccounts)
	f.httpHandler.HandleFunc("POST /api/accounts", f.handleCreateAccount)
	f.httpHandler.HandleFunc("DELETE /api/accounts/{id}", f.handleCloseAccount)
}
