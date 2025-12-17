package manage_transactions

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func (f *Feature) handleRegisterExpense(w http.ResponseWriter, r *http.Request) {
	type RegisterExpenseRequest struct {
		AccountID   uuid.UUID       `json:"account_id"`
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
	}
	var req RegisterExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.RegisterExpense(
		r.Context(),
		uuid.Must(uuid.NewV7()),
		req.AccountID,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleRegisterIncome(w http.ResponseWriter, r *http.Request) {
	type RegisterIncomeRequest struct {
		AccountID   uuid.UUID       `json:"account_id"`
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
	}
	var req RegisterIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.RegisterIncome(
		r.Context(),
		uuid.Must(uuid.NewV7()),
		req.AccountID,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleRegisterTransfer(w http.ResponseWriter, r *http.Request) {
	type RegisterTransferRequest struct {
		FromAccountID uuid.UUID       `json:"from_account_id"`
		FromCurrency  values.Currency `json:"from_currency"`
		FromAmount    decimal.Decimal `json:"from_amount"`
		ToAccountID   uuid.UUID       `json:"to_account_id"`
		ToCurrency    values.Currency `json:"to_currency"`
		ToAmount      decimal.Decimal `json:"to_amount"`
		Category      string          `json:"category"`
		Description   string          `json:"description"`
	}
	var req RegisterTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.RegisterTransfer(
		r.Context(),
		uuid.Must(uuid.NewV7()),
		req.FromAccountID,
		req.FromCurrency,
		req.FromAmount,
		req.ToAccountID,
		req.ToCurrency,
		req.ToAmount,
		req.Category,
		req.Description,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
