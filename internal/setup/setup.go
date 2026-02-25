package setup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/features/import_transactions"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/somatom98/brokeli/internal/features/manage_transactions"
	"github.com/somatom98/brokeli/pkg/event_store"
	"github.com/somatom98/brokeli/pkg/event_store/postgres"
)

type App struct {
	HttpHandler   *http.ServeMux
	httpServer    *http.Server
	transactionES event_store.Store[*transaction.Transaction]
	db            *sql.DB
}

func Setup(ctx context.Context) (*App, error) {
	httpHandler := HttpHandler()

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	accountsRepository, err := accounts.NewPostgresRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts repository: %w", err)
	}

	var transactionES event_store.Store[*transaction.Transaction]

	transactionEventsFactory := map[string]func() interface{}{
		transaction_events.TypeMoneySpent:               func() any { return &transaction_events.MoneySpent{} },
		transaction_events.TypeMoneyReceived:            func() any { return &transaction_events.MoneyReceived{} },
		transaction_events.TypeMoneyTransfered:          func() any { return &transaction_events.MoneyTransfered{} },
		transaction_events.TypeReimbursementReceived:    func() any { return &transaction_events.ReimbursementReceived{} },
		transaction_events.TypeExpectedReimbursementSet: func() any { return &transaction_events.ExpectedReimbursementSet{} },
		transaction_events.TypeMoneyDeposited:           func() any { return &transaction_events.MoneyDeposited{} },
		transaction_events.TypeMoneyWithdrawn:           func() any { return &transaction_events.MoneyWithdrawn{} },
	}

	transactionES, err = postgres.NewPostgresStore(db, transaction.New, transactionEventsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to setup transaction postgres store: %w", err)
	}

	transactionDispatcher := TransactionDispatcher(transactionES)

	accountsProjection := AccountsProjection(ctx, transactionES, accountsRepository)

	manage_transactions.
		New(httpHandler, transactionDispatcher).
		Setup()

	manage_accounts.
		New(httpHandler, accountsProjection).
		Setup()

	import_transactions.
		New(httpHandler, transactionDispatcher).
		Setup()

	return &App{
		HttpHandler:   httpHandler,
		transactionES: transactionES,
		db:            db,
	}, nil
}

func (a *App) Start() <-chan error {
	port := os.Getenv("PORT")

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: a.HttpHandler,
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

	if a.db != nil {
		a.db.Close()
	}

	if closer, ok := a.transactionES.(interface{ Close() error }); ok {
		closer.Close()
	}

	return nil
}
