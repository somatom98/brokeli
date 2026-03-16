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
