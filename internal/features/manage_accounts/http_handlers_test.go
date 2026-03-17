package manage_accounts_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account"
	"github.com/somatom98/brokeli/internal/domain/projections/balance_updates"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/somatom98/brokeli/pkg/event_store"
	"github.com/stretchr/testify/assert"
)

type BalanceUpdatesRepositoryMock struct {
	Balances      []balance_updates.BalancePeriod
	Distributions []balance_updates.AccountDistribution
}

func (m *BalanceUpdatesRepositoryMock) InsertBalanceUpdate(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, userID string, valueDate time.Time, origin string) error {
	return nil
}

func (m *BalanceUpdatesRepositoryMock) GetBalancesByAccount(ctx context.Context, accountID uuid.UUID) ([]balance_updates.BalancePeriod, error) {
	return m.Balances, nil
}

func (m *BalanceUpdatesRepositoryMock) GetAllBalances(ctx context.Context) ([]balance_updates.BalancePeriod, error) {
	return m.Balances, nil
}

func (m *BalanceUpdatesRepositoryMock) GetAccountDistributions(ctx context.Context, accountID uuid.UUID) ([]balance_updates.AccountDistribution, error) {
	return m.Distributions, nil
}

func TestManageAccounts_Handlers(t *testing.T) {
	// arrange
	mux := http.NewServeMux()
	repo := &BalanceUpdatesRepositoryMock{
		Balances: []balance_updates.BalancePeriod{
			{
				Month:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				Currency: values.Currency("EUR"),
				Amount:   decimal.NewFromInt(100),
			},
		},
	}
	dispatcher := &DispatcherMock{}
	transactionES := event_store.NewInMemory[*transaction.Transaction](transaction.New)
	accountES := event_store.NewInMemory[*account.Account](account.New)
	balanceUpdatesProjection := balance_updates.New(transactionES, accountES, repo)
	feature := manage_accounts.New(mux, nil, balanceUpdatesProjection, dispatcher, transactionES)
	feature.Setup(context.Background())

	t.Run("GET /api/balances", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/balances", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var result []balance_updates.BalancePeriod
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "100", result[0].Amount.String())
	})

	t.Run("GET /api/accounts/{id}/balances", func(t *testing.T) {
		id := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/api/accounts/"+id.String()+"/balances", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var result []balance_updates.BalancePeriod
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "100", result[0].Amount.String())
	})

	t.Run("GET /api/accounts/{id}/balances - invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/accounts/invalid-uuid/balances", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("POST /api/accounts", func(t *testing.T) {
		body := `{"name":"test-account", "currency":"EUR"}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var result map[string]string
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)
		assert.NotEmpty(t, result["id"])

		assert.Len(t, dispatcher.Opens, 1)
		assert.Equal(t, "test-account", dispatcher.Opens[0].Name)
		assert.Equal(t, values.Currency("EUR"), dispatcher.Opens[0].Currency)
	})

	t.Run("POST /api/accounts/{id}/deposits", func(t *testing.T) {
		id := uuid.New()
		body := `{"currency":"EUR", "amount":"100.50", "user":"test-user"}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts/"+id.String()+"/deposits", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Len(t, dispatcher.Deposits, 1)
		assert.Equal(t, id, dispatcher.Deposits[0].ID)
		assert.Equal(t, "100.5", dispatcher.Deposits[0].Amount.String())
		assert.Equal(t, "test-user", dispatcher.Deposits[0].User)
	})

	t.Run("POST /api/accounts/{id}/withdrawals", func(t *testing.T) {
		id := uuid.New()
		body := `{"currency":"EUR", "amount":"50.25", "user":"test-user"}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts/"+id.String()+"/withdrawals", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Len(t, dispatcher.Withdrawals, 1)
		assert.Equal(t, id, dispatcher.Withdrawals[0].ID)
		assert.Equal(t, "50.25", dispatcher.Withdrawals[0].Amount.String())
		assert.Equal(t, "test-user", dispatcher.Withdrawals[0].User)
	})
}
