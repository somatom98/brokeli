package manage_transactions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/transactions"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func (f *Feature) handleRegisterExpense(w http.ResponseWriter, r *http.Request) {
	type RegisterExpenseRequest struct {
		AccountID   uuid.UUID       `json:"account_id"`
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
		HappenedAt  time.Time       `json:"happened_at"`
	}
	var req RegisterExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.dispatcher.RegisterExpense(
		r.Context(),
		uuid.Must(uuid.NewV7()),
		req.AccountID,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
		req.HappenedAt,
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
		HappenedAt  time.Time       `json:"happened_at"`
	}
	var req RegisterIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.dispatcher.RegisterIncome(
		r.Context(),
		uuid.Must(uuid.NewV7()),
		req.AccountID,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
		req.HappenedAt,
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
		HappenedAt    time.Time       `json:"happened_at"`
	}
	var req RegisterTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

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
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleRegisterReimbursement(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("transaction_id")
	if transactionID == "" {
		http.Error(w, "bad request: missing transaction id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(transactionID)
	if err != nil {
		http.Error(w, "bad request: invalid transaction id", http.StatusBadRequest)
		return
	}

	type RegisterReimbursementRequest struct {
		AccountID   uuid.UUID       `json:"account_id"`
		From        string          `json:"from"`
		Currency    values.Currency `json:"currency"`
		Amount      decimal.Decimal `json:"amount"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
		HappenedAt  time.Time       `json:"happened_at"`
	}

	var req RegisterReimbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.dispatcher.RegisterReimbursement(
		r.Context(),
		id,
		req.AccountID,
		req.From,
		req.Currency,
		req.Amount,
		req.Category,
		req.Description,
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleSetExpectedReimbursement(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("transaction_id")
	if transactionID == "" {
		http.Error(w, "bad request: missing transaction id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(transactionID)
	if err != nil {
		http.Error(w, "bad request: invalid transaction id", http.StatusBadRequest)
		return
	}

	type SetExpectedReimbursementRequest struct {
		AccountID  uuid.UUID       `json:"account_id"`
		Currency   values.Currency `json:"currency"`
		Amount     decimal.Decimal `json:"amount"`
		HappenedAt time.Time       `json:"happened_at"`
	}

	var req SetExpectedReimbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.dispatcher.SetExpectedReimbursement(
		r.Context(),
		id,
		req.AccountID,
		req.Currency,
		req.Amount,
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleRegisterInvestment(w http.ResponseWriter, r *http.Request) {
	type RegisterInvestmentRequest struct {
		AccountID     uuid.UUID       `json:"account_id"`
		Ticker        string          `json:"ticker"`
		Units         decimal.Decimal `json:"units"`
		Price         decimal.Decimal `json:"price"`
		PriceCurrency values.Currency `json:"price_currency"`
		Fee           decimal.Decimal `json:"fee"`
		FeeCurrency   values.Currency `json:"fee_currency"`
		HappenedAt    time.Time       `json:"happened_at"`
	}

	var req RegisterInvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.HappenedAt.IsZero() {
		req.HappenedAt = time.Now()
	}

	if err := f.dispatcher.RegisterInvestment(
		r.Context(),
		uuid.Must(uuid.NewV7()),
		req.AccountID,
		req.Ticker,
		req.Units,
		req.Price,
		req.PriceCurrency,
		req.Fee,
		req.FeeCurrency,
		req.HappenedAt,
	); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (f *Feature) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := transactions.ListTransactionsParams{}

	if startStr := query.Get("start_date"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			params.StartDate = &t
		}
	}

	if endStr := query.Get("end_date"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			params.EndDate = &t
		}
	}

	if accounts := query["account_id"]; len(accounts) > 0 {
		for _, a := range accounts {
			if id, err := uuid.Parse(a); err == nil {
				params.AccountIDs = append(params.AccountIDs, id)
			}
		}
	}

	if tType := query.Get("transaction_type"); tType != "" {
		params.TransactionType = &tType
	}

	isPaginated := query.Get("paginated") == "true"
	if isPaginated {
		page := 1
		pageSize := 50

		if p, err := strconv.Atoi(query.Get("page")); err == nil && p > 0 {
			page = p
		}
		if ps, err := strconv.Atoi(query.Get("page_size")); err == nil && ps > 0 {
			pageSize = ps
		}

		paginatedParams := transactions.ListTransactionsPaginatedParams{
			ListTransactionsParams: params,
			Limit:                  int32(pageSize),
			Offset:                 int32((page - 1) * pageSize),
		}

		results, err := f.transactionsView.ListTransactionsPaginated(r.Context(), paginatedParams)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		return
	}

	transactions, err := f.transactionsView.ListTransactions(r.Context(), params)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
