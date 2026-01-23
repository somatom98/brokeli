package manage_accounts

import (
	"encoding/json"
	"net/http"
	"time"

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

func (f *Feature) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	id := uuid.Must(uuid.NewV7())

	if err := f.dispatcher.CreateAccount(r.Context(), id, time.Now()); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(struct {
		ID uuid.UUID `json:"id"`
	}{
		ID: id,
	})
}

func (f *Feature) handleCloseAccount(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("id")
	if transactionID == "" {
		http.Error(w, "bad request: missing account id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(transactionID)
	if err != nil {
		http.Error(w, "bad request: invalid account id", http.StatusBadRequest)
		return
	}

	if err := f.dispatcher.CloseAccount(
		r.Context(),
		id,
		time.Now(),
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
