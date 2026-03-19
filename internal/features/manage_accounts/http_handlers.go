package manage_accounts

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
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

	balances, err := f.balanceUpdatesView.GetBalancesByAccount(r.Context(), id)
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

func (f *Feature) handleGetAccountDistributions(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	distributions, err := f.balanceUpdatesView.GetAccountDistributions(r.Context(), id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	jsonDistributions, err := json.Marshal(distributions)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDistributions)
}

func (f *Feature) handleGetAllBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := f.balanceUpdatesView.GetAllBalances(r.Context())
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

func (f *Feature) handleDeposit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	type DepositRequest struct {
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
		User        string          `json:"user"`
		HappenedAt  time.Time       `json:"happened_at"`
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.accountDispatcher.Deposit(
		r.Context(),
		id,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
		req.User,
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleWithdrawal(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	type WithdrawalRequest struct {
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
		User        string          `json:"user"`
		HappenedAt  time.Time       `json:"happened_at"`
	}

	var req WithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.accountDispatcher.Withdraw(
		r.Context(),
		id,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
		req.User,
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleOpenAccount(w http.ResponseWriter, r *http.Request) {
	type OpenAccountRequest struct {
		ID         uuid.UUID       `json:"id"`
		Name       string          `json:"name"`
		Currency   values.Currency `json:"currency"`
		HappenedAt time.Time       `json:"happened_at"`
	}

	var req OpenAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.accountDispatcher.Open(
		r.Context(),
		req.ID,
		req.Name,
		req.Currency,
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": req.ID.String()})
}
