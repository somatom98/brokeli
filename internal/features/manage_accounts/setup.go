package manage_accounts

import (
	"net/http"

	"github.com/somatom98/brokeli/internal/domain/views/accounts"
)

type Feature struct {
	httpHandler  *http.ServeMux
	accountsView *accounts.View
}

func New(
	httpHandler *http.ServeMux,
	accountsView *accounts.View,
) *Feature {
	return &Feature{
		httpHandler:  httpHandler,
		accountsView: accountsView,
	}
}

func (f *Feature) Setup() {
	f.httpHandler.HandleFunc("GET /api/accounts", f.handleGetAccounts)
}
