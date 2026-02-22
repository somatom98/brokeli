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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/internal/setup"
)

func TestIntegration_ManageAccounts(t *testing.T) {
	// Arrange
	os.Setenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/brokeli?sslmode=disable")

	ctx := context.Background()
	app, err := setup.Setup(ctx)
	require.NoError(t, err, "Failed to setup app")

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
