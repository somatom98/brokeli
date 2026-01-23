package accounts

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) (*PostgresRepository, error) {
	query := `
	CREATE TABLE IF NOT EXISTS accounts_projection (
		id UUID PRIMARY KEY,
		balance JSONB NOT NULL DEFAULT '{}',
		created_at TIMESTAMP,
		closed_at TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts_projection table: %w", err)
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error {
	// We use ON CONFLICT because we might receive multiple events or out of order, 
	// though creation should be first. 
	// If it exists, we just update created_at if it was null (which shouldn't happen with strict ordering but good for idempotency).
	query := `
	INSERT INTO accounts_projection (id, created_at, balance)
	VALUES ($1, $2, '{}')
	ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at
	`
	_, err := r.db.ExecContext(ctx, query, id, createdAt)
	return err
}

func (r *PostgresRepository) CloseAccount(ctx context.Context, id uuid.UUID, closedAt time.Time) error {
	query := `
	UPDATE accounts_projection
	SET closed_at = $2
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, closedAt)
	return err
}

func (r *PostgresRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error {
	// This requires a read-modify-write transaction to ensure consistency.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Read current balance
	var balanceJSON []byte
	query := `SELECT balance FROM accounts_projection WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, id).Scan(&balanceJSON)
	if err == sql.ErrNoRows {
		// If account doesn't exist yet (out of order event?), strictly we should probably create it or error.
		// For robustness, let's create a placeholder.
		balanceJSON = []byte("{}")
		_, err = tx.ExecContext(ctx, `INSERT INTO accounts_projection (id, balance) VALUES ($1, $2)`, id, balanceJSON)
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
	updateQuery := `UPDATE accounts_projection SET balance = $2 WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateQuery, id, newBalanceJSON)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetAll(ctx context.Context) (map[uuid.UUID]Account, error) {
	query := `SELECT id, balance, created_at, closed_at FROM accounts_projection`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make(map[uuid.UUID]Account)
	for rows.Next() {
		var id uuid.UUID
		var balanceJSON []byte
		var createdAt sql.NullTime
		var closedAt sql.NullTime

		if err := rows.Scan(&id, &balanceJSON, &createdAt, &closedAt); err != nil {
			return nil, err
		}

		var balance map[values.Currency]decimal.Decimal
		if err := json.Unmarshal(balanceJSON, &balance); err != nil {
			return nil, err
		}

		acc := Account{
			Balance: balance,
		}
		if createdAt.Valid {
			t := createdAt.Time
			acc.CreatedAt = &t
		}
		if closedAt.Valid {
			t := closedAt.Time
			acc.ClosedAt = &t
		}

		accounts[id] = acc
	}
	return accounts, nil
}
