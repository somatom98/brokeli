package manage_accounts

import (
	"net/http"

	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
)

type Feature struct {
	httpHandler  *http.ServeMux
	accountsView *accounts.Projection
}

func New(
	httpHandler *http.ServeMux,
	accountsView *accounts.Projection,
) *Feature {
	return &Feature{
		httpHandler:  httpHandler,
		accountsView: accountsView,
	}
}

func (f *Feature) Setup() {
	f.httpHandler.HandleFunc("GET /api/accounts", f.handleGetAccounts)
}
