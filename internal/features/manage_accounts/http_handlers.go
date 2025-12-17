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
	w.Write(jsonAccounts)
	w.WriteHeader(http.StatusCreated)
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
	type CloseRequest struct {
		AccountID uuid.UUID `json:"account_id"`
	}

	var req CloseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.CloseAccount(
		r.Context(),
		req.AccountID,
		time.Now(),
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
