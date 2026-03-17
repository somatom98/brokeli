package manage_budgets

import (
	"context"
	"net/http"

	"github.com/somatom98/brokeli/internal/domain/budget"
	"github.com/somatom98/brokeli/internal/domain/projections/transactions"
)

type Feature struct {
	httpHandler      *http.ServeMux
	budgetRepository budget.Repository
	transactionsView *transactions.Projection
}

func New(
	httpHandler *http.ServeMux,
	budgetRepository budget.Repository,
	transactionsView *transactions.Projection,
) *Feature {
	f := &Feature{
		httpHandler:      httpHandler,
		budgetRepository: budgetRepository,
		transactionsView: transactionsView,
	}

	return f
}

func (f *Feature) Setup(ctx context.Context) {
	f.httpHandler.HandleFunc("GET /api/budgets", f.handleGetBudgets)
	f.httpHandler.HandleFunc("GET /api/budgets/categories", f.handleGetCategories)
	f.httpHandler.HandleFunc("POST /api/budgets", f.handleSaveBudget)
	f.httpHandler.HandleFunc("DELETE /api/budgets/{id}", f.handleDeleteBudget)
}
