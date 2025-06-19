package create_transactions

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func (f *Feature) handleCreateExpense(w http.ResponseWriter, r *http.Request) {
	type CreateExpenseRequest struct {
		AccountID   uuid.UUID       `json:"account_id"`
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
	}
	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.CreateExpense(r.Context(), uuid.Must(uuid.NewV7()), transaction.CreateExpense{
		AccountID:   req.AccountID,
		Currency:    req.Currency,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleCreateIncome(w http.ResponseWriter, r *http.Request) {
	type CreateIncomeRequest struct {
		AccountID   uuid.UUID       `json:"account_id"`
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
	}
	var req CreateIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.CreateIncome(r.Context(), uuid.Must(uuid.NewV7()), transaction.CreateIncome{
		AccountID:   req.AccountID,
		Currency:    req.Currency,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	type CreateTransferRequest struct {
		FromAccountID uuid.UUID       `json:"from_account_id"`
		FromCurrency  values.Currency `json:"from_currency"`
		FromAmount    decimal.Decimal `json:"from_amount"`
		ToAccountID   uuid.UUID       `json:"to_account_id"`
		ToCurrency    values.Currency `json:"to_currency"`
		ToAmount      decimal.Decimal `json:"to_amount"`
		Category      string          `json:"category"`
		Description   string          `json:"description"`
	}
	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := f.dispatcher.CreateTransfer(r.Context(), uuid.Must(uuid.NewV7()), transaction.CreateTransfer{
		FromAccountID: req.FromAccountID,
		FromCurrency:  req.FromCurrency,
		FromAmount:    req.FromAmount,
		ToAccountID:   req.ToAccountID,
		ToCurrency:    req.ToCurrency,
		ToAmount:      req.ToAmount,
		Category:      req.Category,
		Description:   req.Description,
	}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
