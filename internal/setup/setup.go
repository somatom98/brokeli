package setup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/somatom98/brokeli/internal/domain/account"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
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
	accountES     event_store.Store[*account.Account]
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

	var accountES event_store.Store[*account.Account]
	var transactionES event_store.Store[*transaction.Transaction]

	accountEventsFactory := map[string]func() interface{}{
		account_events.Type_Created:        func() interface{} { return &account_events.Created{} },
		account_events.Type_MoneyDeposited: func() interface{} { return &account_events.MoneyDeposited{} },
		account_events.Type_AccountClosed:  func() interface{} { return &account_events.AccountClosed{} },
	}

	accountES, err = postgres.NewPostgresStore(db, account.New, accountEventsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to setup account postgres store: %w", err)
	}

	transactionEventsFactory := map[string]func() interface{}{
		transaction_events.Type_MoneySpent:               func() interface{} { return &transaction_events.MoneySpent{} },
		transaction_events.Type_MoneyReceived:            func() interface{} { return &transaction_events.MoneyReceived{} },
		transaction_events.Type_MoneyTransfered:          func() interface{} { return &transaction_events.MoneyTransfered{} },
		transaction_events.Type_ReimbursementReceived:    func() interface{} { return &transaction_events.ReimbursementReceived{} },
		transaction_events.Type_ExpectedReimbursementSet: func() interface{} { return &transaction_events.ExpectedReimbursementSet{} },
	}

	transactionES, err = postgres.NewPostgresStore(db, transaction.New, transactionEventsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to setup transaction postgres store: %w", err)
	}

	accountDispatcher := AccountDispatcher(accountES)
	transactionDispatcher := TransactionDispatcher(transactionES)

	accountsProjection := AccountsProjection(ctx, transactionES, accountES, accountsRepository)

	manage_transactions.
		New(httpHandler, transactionDispatcher).
		Setup()

	manage_accounts.
		New(httpHandler, accountsProjection, accountDispatcher).
		Setup()

	import_transactions.
		New(httpHandler, transactionDispatcher).
		Setup()

	return &App{
		HttpHandler:   httpHandler,
		accountES:     accountES,
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

	if closer, ok := a.accountES.(interface{ Close() error }); ok {
		closer.Close()
	}
	if closer, ok := a.transactionES.(interface{ Close() error }); ok {
		closer.Close()
	}

	return nil
}
