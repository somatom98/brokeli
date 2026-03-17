package transactions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/db"
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

func (r *PostgresRepository) CreateTransaction(ctx context.Context, tx TransactionRecord) error {
	return r.queries.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:              tx.ID,
		AccountID:       tx.AccountID,
		TransactionType: tx.TransactionType,
		Amount:          tx.Amount.String(),
		Currency:        string(tx.Currency),
		Category:        tx.Category,
		Description:     tx.Description,
		HappenedAt:      tx.HappenedAt,
	})
}

func (r *PostgresRepository) ListTransactionsByAccount(ctx context.Context, accountID uuid.UUID) ([]TransactionRecord, error) {
	rows, err := r.queries.ListTransactionsByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	transactions := make([]TransactionRecord, len(rows))
	for i, row := range rows {
		amount, _ := decimal.NewFromString(row.Amount)
		
		var rate decimal.Decimal
		if row.SystemTotalRate != nil {
			rate, _ = decimal.NewFromString(fmt.Sprintf("%v", row.SystemTotalRate))
		}

		transactions[i] = TransactionRecord{
			ID:              row.ID,
			AccountID:       row.AccountID,
			TransactionType: row.TransactionType,
			Amount:          amount,
			Currency:        values.Currency(row.Currency),
			Category:        row.Category,
			Description:     row.Description,
			HappenedAt:      row.HappenedAt,
			SystemTotalRate: rate,
		}
	}

	return transactions, nil
}

func (r *PostgresRepository) ListTransactions(ctx context.Context) ([]TransactionRecord, error) {
	rows, err := r.queries.ListTransactions(ctx)
	if err != nil {
		return nil, err
	}

	transactions := make([]TransactionRecord, len(rows))
	for i, row := range rows {
		amount, _ := decimal.NewFromString(row.Amount)
		
		var rate decimal.Decimal
		if row.SystemTotalRate != nil {
			rate, _ = decimal.NewFromString(fmt.Sprintf("%v", row.SystemTotalRate))
		}

		transactions[i] = TransactionRecord{
			ID:              row.ID,
			AccountID:       row.AccountID,
			TransactionType: row.TransactionType,
			Amount:          amount,
			Currency:        values.Currency(row.Currency),
			Category:        row.Category,
			Description:     row.Description,
			HappenedAt:      row.HappenedAt,
			SystemTotalRate: rate,
		}
	}

	return transactions, nil
}
