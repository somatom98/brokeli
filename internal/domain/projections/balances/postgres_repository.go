package balances

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/projections/balances/db"
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

func (r *PostgresRepository) InsertBalanceUpdate(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, valueDate time.Time) error {
	return r.queries.InsertBalanceUpdate(ctx, db.InsertBalanceUpdateParams{
		ID:        id,
		AccountID: accountID,
		Currency:  string(currency),
		Amount:    amount.String(),
		ValueDate: valueDate,
	})
}

func (r *PostgresRepository) GetBalancesByAccount(ctx context.Context, accountID uuid.UUID) ([]BalancePeriod, error) {
	rows, err := r.queries.GetBalancesByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	cumulativeBalances := make(map[values.Currency]decimal.Decimal)
	balances := make([]BalancePeriod, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		row := rows[i]
		currency := values.Currency(row.Currency)
		amount, _ := decimal.NewFromString(row.Amount)

		currentBalance := cumulativeBalances[currency].Add(amount)
		cumulativeBalances[currency] = currentBalance

		balances[i] = BalancePeriod{
			Month:    row.Month,
			Currency: currency,
			Amount:   currentBalance,
		}
	}

	return balances, nil
}

func (r *PostgresRepository) GetAllBalances(ctx context.Context) ([]BalancePeriod, error) {
	rows, err := r.queries.GetAllBalances(ctx)
	if err != nil {
		return nil, err
	}

	cumulativeBalances := make(map[values.Currency]decimal.Decimal)
	balances := make([]BalancePeriod, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		row := rows[i]
		currency := values.Currency(row.Currency)
		amount, _ := decimal.NewFromString(row.Amount)

		currentBalance := cumulativeBalances[currency].Add(amount)
		cumulativeBalances[currency] = currentBalance

		balances[i] = BalancePeriod{
			Month:    row.Month,
			Currency: currency,
			Amount:   currentBalance,
		}
	}

	return balances, nil
}
