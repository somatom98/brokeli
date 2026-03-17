package setup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	projections_db "github.com/somatom98/brokeli/internal/db"
	"github.com/somatom98/brokeli/internal/domain/account"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/projections/accounts"
	"github.com/somatom98/brokeli/internal/domain/projections/balance_updates"
	"github.com/somatom98/brokeli/internal/domain/projections/transactions"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/features/import_transactions"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/somatom98/brokeli/internal/features/manage_transactions"
	"github.com/somatom98/brokeli/pkg/database"
	"github.com/somatom98/brokeli/pkg/event_store"
	"github.com/somatom98/brokeli/pkg/event_store/postgres"
	event_store_db "github.com/somatom98/brokeli/pkg/event_store/postgres/db"
)

type App struct {
	HttpHandler   *http.ServeMux
	httpServer    *http.Server
	transactionES event_store.Store[*transaction.Transaction]
	accountES     event_store.Store[*account.Account]
	db            *sql.DB
	cancelRelays  context.CancelFunc
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

	// Run migrations
	if err := database.Migrate(db, event_store_db.MigrationsFS(), "event_store_migrations"); err != nil {
		return nil, fmt.Errorf("failed to run event store migrations: %w", err)
	}
	if err := database.Migrate(db, projections_db.MigrationsFS(), "projections_migrations"); err != nil {
		return nil, fmt.Errorf("failed to run projections migrations: %w", err)
	}

	accountsRepository, err := accounts.NewPostgresRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts repository: %w", err)
	}

	balanceUpdatesRepository, err := balance_updates.NewPostgresRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create balance updates repository: %w", err)
	}

	transactionsRepository, err := transactions.NewPostgresRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions repository: %w", err)
	}

	var transactionES event_store.Store[*transaction.Transaction]

	transactionEventsFactory := map[string]func() any{
		transaction_events.TypeMoneySpent:               func() any { return &transaction_events.MoneySpent{} },
		transaction_events.TypeMoneyReceived:            func() any { return &transaction_events.MoneyReceived{} },
		transaction_events.TypeMoneyTransfered:          func() any { return &transaction_events.MoneyTransfered{} },
		transaction_events.TypeReimbursementReceived:    func() any { return &transaction_events.ReimbursementReceived{} },
		transaction_events.TypeExpectedReimbursementSet: func() any { return &transaction_events.ExpectedReimbursementSet{} },
	}

	transactionES, err = postgres.NewPostgresStore(db, transaction.New, transactionEventsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to setup transaction postgres store: %w", err)
	}

	var accountES event_store.Store[*account.Account]

	accountEventsFactory := map[string]func() any{
		account_events.TypeOpened:         func() any { return &account_events.Opened{} },
		account_events.TypeNameUpdated:    func() any { return &account_events.NameUpdated{} },
		account_events.TypeMoneyDeposited: func() any { return &account_events.MoneyDeposited{} },
		account_events.TypeMoneyWithdrawn: func() any { return &account_events.MoneyWithdrawn{} },
	}

	accountES, err = postgres.NewPostgresStore(db, account.New, accountEventsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to setup account postgres store: %w", err)
	}

	transactionDispatcher := TransactionDispatcher(transactionES)
	accountDispatcher := AccountDispatcher(accountES)

	accountsProjection := AccountsProjection(ctx, transactionES, accountES, accountsRepository)
	balanceUpdatesProjection := BalanceUpdatesProjection(ctx, transactionES, accountES, balanceUpdatesRepository)
	transactionsProjection := TransactionsProjection(ctx, transactionES, accountES, transactionsRepository)

	manage_transactions.
		New(httpHandler, transactionDispatcher, transactionsProjection).
		Setup()

	manage_accounts.
		New(httpHandler, accountsProjection, balanceUpdatesProjection, accountDispatcher, transactionES).
		Setup(ctx)

	import_transactions.
		New(httpHandler, transactionDispatcher, accountDispatcher).
		Setup()

	return &App{
		HttpHandler:   httpHandler,
		transactionES: transactionES,
		accountES:     accountES,
		db:            db,
		cancelRelays:  func() {},
	}, nil
}

func (a *App) Start() <-chan error {
	port := os.Getenv("PORT")

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: a.HttpHandler,
	}

	errCh := make(chan error)

	relayCtx, cancel := context.WithCancel(context.Background())
	a.cancelRelays = cancel

	// Start relay workers
	if es, ok := a.transactionES.(*postgres.PostgresStore[*transaction.Transaction]); ok {
		go func() {
			if err := es.RunRelay(relayCtx); err != nil && err != context.Canceled {
				log.Printf("Transaction Relay error: %v", err)
			}
		}()
	}

	if es, ok := a.accountES.(*postgres.PostgresStore[*account.Account]); ok {
		go func() {
			if err := es.RunRelay(relayCtx); err != nil && err != context.Canceled {
				log.Printf("Account Relay error: %v", err)
			}
		}()
	}

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
	if a.cancelRelays != nil {
		a.cancelRelays()
	}

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

	if closer, ok := a.accountES.(interface{ Close() error }); ok {
		closer.Close()
	}

	return nil
}
