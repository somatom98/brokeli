package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/internal/setup"
)

func TestIntegration_ManageAccounts(t *testing.T) {
	// Arrange
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/brokeli?sslmode=disable"
	}
	os.Setenv("DB_DSN", dsn)

	ctx := context.Background()
	app, err := setup.Setup(ctx)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}

	server := httptest.NewServer(app.HttpHandler)
	defer server.Close()

	client := server.Client()

	t.Run("Get Accounts", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)

		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/accounts", nil)
		require.NoError(t, err)

		// Act
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var accounts map[uuid.UUID]interface{}
		err = json.NewDecoder(resp.Body).Decode(&accounts)
		require.NoError(t, err)
	})
}

func TestIntegration_ImportTransactions(t *testing.T) {
	// Arrange
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/brokeli?sslmode=disable"
	}
	os.Setenv("DB_DSN", dsn)

	ctx := context.Background()
	app, err := setup.Setup(ctx)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	app.Start()

	defer app.Stop(ctx)

	t.Run("Import Transactions", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/import-transactions", nil)
		require.NoError(t, err)

		// Act
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Wait for processing
		time.Sleep(5 * time.Second)

		// Verify balances
		req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/api/accounts", nil)
		require.NoError(t, err)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var accountsResponse map[uuid.UUID]struct {
			Balance map[string]decimal.Decimal `json:"balance"`
		}
		err = json.NewDecoder(resp.Body).Decode(&accountsResponse)
		require.NoError(t, err)

		// Check one of the accounts balance
		lunarID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Lunar"))
		assert.NotEmpty(t, accountsResponse[lunarID].Balance["DKK"])
		assert.True(t, accountsResponse[lunarID].Balance["DKK"].Equal(decimal.NewFromFloat(101132.95)))
	})
}
