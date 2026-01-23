package setup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/somatom98/brokeli/internal/domain/account"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/somatom98/brokeli/internal/features/manage_transactions"
	"github.com/somatom98/brokeli/pkg/event_store/postgres"
)

type App struct {
	httpHandler   *http.ServeMux
	httpServer    *http.Server
	accountES     *postgres.PostgresStore[*account.Account]
	transactionES *postgres.PostgresStore[*transaction.Transaction]
}

func Setup(ctx context.Context) (*App, error) {
	httpHandler := HttpHandler()

	accountES, err := postgres.Setup(os.Getenv("DB_DSN"), account.New)
	if err != nil {
		return nil, fmt.Errorf("failed to setup account event store: %w", err)
	}

	transactionES, err := postgres.Setup(os.Getenv("DB_DSN"), transaction.New)
	if err != nil {
		return nil, fmt.Errorf("failed to setup transaction event store: %w", err)
	}

	accountDispatcher := AccountDispatcher(accountES)
	transactionDispatcher := TransactionDispatcher(transactionES)

	accountsProjection := AccountsProjection(ctx, transactionES, accountES)

	manage_transactions.
		New(httpHandler, transactionDispatcher).
		Setup()

	manage_accounts.
		New(httpHandler, accountsProjection, accountDispatcher).
		Setup()

	return &App{
		httpHandler:   httpHandler,
		accountES:     accountES,
		transactionES: transactionES,
	}, nil
}

func (a *App) Start() <-chan error {
	port := os.Getenv("PORT")

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: a.httpHandler,
	}

	errCh := make(chan error)

	go func() {
		defer close(errCh)

		log.Printf("Starting server on :%s", port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("ListenAndServe: %w", err)
			return
		}
	}()

	return errCh
}

func (a *App) Stop(ctx context.Context) error {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	a.accountES.Close()
	a.transactionES.Close()

	return nil
}
