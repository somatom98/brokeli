package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/somatom98/brokeli/internal/setup"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	t        *testing.T
	ctx      context.Context
	client   *http.Client
	baseURL  string
	accounts map[string]string // alias -> id
}

func NewSuite(t *testing.T) *Suite {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	container, err := SetupPostgres(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { container.Close(ctx) })

	os.Setenv("DB_DSN", container.DSN)
	os.Setenv("PORT", "8082")

	app, err := setup.Setup(ctx)
	require.NoError(t, err)

	errCh := app.Start()
	select {
	case err := <-errCh:
		t.Fatalf("App failed to start: %v", err)
	case <-time.After(2 * time.Second):
		// App started
	}
	t.Cleanup(func() { app.Stop(ctx) })

	return &Suite{
		t:        t,
		ctx:      ctx,
		client:   &http.Client{Timeout: 10 * time.Second},
		baseURL:  "http://localhost:8082/api",
		accounts: make(map[string]string),
	}
}

func (s *Suite) Given() *Given {
	return &Given{s: s}
}

func (s *Suite) When() *When {
	return &When{s: s}
}

func (s *Suite) Then() *Then {
	return &Then{s: s}
}

const (
	EUR = "EUR"
	USD = "USD"
	DKK = "DKK"
)
