package accounts

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/accounts/db"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type PostgresRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewPostgresRepository(dbConn *sql.DB) (*PostgresRepository, error) {
	return &PostgresRepository{
		db:      dbConn,
		queries: db.New(dbConn),
	}, nil
}

func (r *PostgresRepository) CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error {
	return r.queries.CreateAccount(ctx, db.CreateAccountParams{
		ID:        id,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	})
}

func (r *PostgresRepository) CloseAccount(ctx context.Context, id uuid.UUID, closedAt time.Time) error {
	return r.queries.CloseAccount(ctx, db.CloseAccountParams{
		ID:       id,
		ClosedAt: sql.NullTime{Time: closedAt, Valid: true},
	})
}

func (r *PostgresRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error {
	// This requires a read-modify-write transaction to ensure consistency.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	// Read current balance
	balanceJSON, err := qtx.GetAccountBalanceForUpdate(ctx, id)
	if err == sql.ErrNoRows {
		// If account doesn't exist yet (out of order event?), strictly we should probably create it or error.
		// For robustness, let's create a placeholder.
		balanceJSON = []byte("{}")
		err = qtx.UpsertPlaceholderAccount(ctx, db.UpsertPlaceholderAccountParams{
			ID:      id,
			Balance: balanceJSON,
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var balance map[values.Currency]decimal.Decimal
	if err := json.Unmarshal(balanceJSON, &balance); err != nil {
		return err
	}
	if balance == nil {
		balance = make(map[values.Currency]decimal.Decimal)
	}

	// Update balance
	if _, ok := balance[currency]; !ok {
		balance[currency] = decimal.Zero
	}
	balance[currency] = balance[currency].Add(amount)

	newBalanceJSON, err := json.Marshal(balance)
	if err != nil {
		return err
	}

	// Write back
	err = qtx.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		ID:      id,
		Balance: newBalanceJSON,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetAll(ctx context.Context) (map[uuid.UUID]Account, error) {
	rows, err := r.queries.GetAllAccounts(ctx)
	if err != nil {
		return nil, err
	}

	accounts := make(map[uuid.UUID]Account)
	for _, row := range rows {
		var balance map[values.Currency]decimal.Decimal
		if err := json.Unmarshal(row.Balance, &balance); err != nil {
			return nil, err
		}

		acc := Account{
			Balance: balance,
		}
		if row.CreatedAt.Valid {
			t := row.CreatedAt.Time
			acc.CreatedAt = &t
		}
		if row.ClosedAt.Valid {
			t := row.ClosedAt.Time
			acc.ClosedAt = &t
		}

		accounts[row.ID] = acc
	}
	return accounts, nil
}
