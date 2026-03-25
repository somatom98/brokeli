package transactions

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *PostgresRepository) ListTransactions(ctx context.Context, params ListTransactionsParams) ([]TransactionRecord, error) {
	arg := db.ListTransactionsParams{
		AccountIds: params.AccountIDs,
	}

	if params.StartDate != nil {
		arg.StartDate = sql.NullTime{
			Time:  *params.StartDate,
			Valid: true,
		}
	}

	if params.EndDate != nil {
		arg.EndDate = sql.NullTime{
			Time:  *params.EndDate,
			Valid: true,
		}
	}

	if params.TransactionType != nil {
		arg.TransactionType = sql.NullString{
			String: *params.TransactionType,
			Valid:  true,
		}
	}

	rows, err := r.queries.ListTransactions(ctx, arg)
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

func (r *PostgresRepository) ListTransactionsPaginated(ctx context.Context, params ListTransactionsPaginatedParams) (PaginatedTransactions, error) {
	arg := db.ListTransactionsPaginatedParams{
		AccountIds: params.AccountIDs,
		LimitVal:   params.Limit,
		OffsetVal:  params.Offset,
	}

	if params.StartDate != nil {
		arg.StartDate = sql.NullTime{
			Time:  *params.StartDate,
			Valid: true,
		}
	}

	if params.EndDate != nil {
		arg.EndDate = sql.NullTime{
			Time:  *params.EndDate,
			Valid: true,
		}
	}

	if params.TransactionType != nil {
		arg.TransactionType = sql.NullString{
			String: *params.TransactionType,
			Valid:  true,
		}
	}

	rows, err := r.queries.ListTransactionsPaginated(ctx, arg)
	if err != nil {
		return PaginatedTransactions{}, err
	}

	transactions := make([]TransactionRecord, len(rows))
	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}

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

	return PaginatedTransactions{
		Transactions: transactions,
		TotalCount:   totalCount,
	}, nil
}

func (r *PostgresRepository) ListCategories(ctx context.Context) ([]string, error) {
	return r.queries.ListCategories(ctx)
}
