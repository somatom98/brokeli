package setup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/somatom98/brokeli/internal/features/manage_transactions"
	event_store "github.com/somatom98/brokeli/pkg/event_store/sqlite"
)

type App struct {
	httpHandler   *http.ServeMux
	httpServer    *http.Server
	transactionES *event_store.SQLiteStore[*transaction.Transaction]
}

func Setup(ctx context.Context) (*App, error) {
	httpHandler := HttpHandler()

	transactionES, err := event_store.Setup(os.Getenv("DB_PATH"), transaction.New)
	if err != nil {
		return nil, fmt.Errorf("failed to setup transaction event store: %w", err)
	}

	transactionDispatcher := TransactionDispatcher(transactionES)

	accountsProjection := AccountsProjection(ctx, transactionES)

	manage_transactions.
		New(httpHandler, transactionDispatcher).
		Setup()

	manage_accounts.
		New(httpHandler, accountsProjection).
		Setup()

	return &App{
		httpHandler:   httpHandler,
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

	a.transactionES.Close()

	return nil
}
