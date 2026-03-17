package manage_budgets

import (
	"context"
	"net/http"

	"github.com/somatom98/brokeli/internal/domain/budget"
)

type Feature struct {
	httpHandler      *http.ServeMux
	budgetRepository budget.Repository
}

func New(
	httpHandler *http.ServeMux,
	budgetRepository budget.Repository,
) *Feature {
	f := &Feature{
		httpHandler:      httpHandler,
		budgetRepository: budgetRepository,
	}

	return f
}

func (f *Feature) Setup(ctx context.Context) {
	f.httpHandler.HandleFunc("GET /api/budgets", f.handleGetBudgets)
	f.httpHandler.HandleFunc("POST /api/budgets", f.handleSaveBudget)
	f.httpHandler.HandleFunc("DELETE /api/budgets/{id}", f.handleDeleteBudget)
}
