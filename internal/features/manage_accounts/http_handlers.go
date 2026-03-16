package manage_accounts

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (f *Feature) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := f.accountsView.GetAll(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	jsonAccounts, err := json.Marshal(accounts)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonAccounts)
}

func (f *Feature) handleGetAccountBalances(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	balances, err := f.balancesView.GetBalancesByAccount(r.Context(), id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	jsonBalances, err := json.Marshal(balances)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBalances)
}

func (f *Feature) handleGetAllBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := f.balancesView.GetAllBalances(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	jsonBalances, err := json.Marshal(balances)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBalances)
}
