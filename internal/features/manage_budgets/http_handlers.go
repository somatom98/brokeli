package manage_budgets

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/budget"
)

func (f *Feature) handleGetBudgets(w http.ResponseWriter, r *http.Request) {
	budgets, err := f.budgetRepository.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func (f *Feature) handleSaveBudget(w http.ResponseWriter, r *http.Request) {
	var b budget.Budget
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	err := f.budgetRepository.Save(r.Context(), b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleDeleteBudget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = f.budgetRepository.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
