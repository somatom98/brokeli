package manage_accounts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/balances"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/stretchr/testify/assert"
)

type BalancesRepositoryMock struct {
	Balances []balances.BalancePeriod
}

func (m *BalancesRepositoryMock) InsertBalanceUpdate(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, userID string, valueDate time.Time) error {
	return nil
}

func (m *BalancesRepositoryMock) GetBalancesByAccount(ctx context.Context, accountID uuid.UUID) ([]balances.BalancePeriod, error) {
	return m.Balances, nil
}

func (m *BalancesRepositoryMock) GetAllBalances(ctx context.Context) ([]balances.BalancePeriod, error) {
	return m.Balances, nil
}

func TestManageAccounts_GetBalancesHandlers(t *testing.T) {
	// arrange
	mux := http.NewServeMux()
	repo := &BalancesRepositoryMock{
		Balances: []balances.BalancePeriod{
			{
				Month:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				Currency: values.Currency("EUR"),
				Amount:   decimal.NewFromInt(100),
			},
		},
	}
	balancesProjection := balances.New(nil, nil, repo)
	feature := manage_accounts.New(mux, nil, balancesProjection, nil, nil)
	feature.Setup(context.Background())

	t.Run("GET /api/balances", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/balances", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var result []balances.BalancePeriod
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

		var result []balances.BalancePeriod
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
}
