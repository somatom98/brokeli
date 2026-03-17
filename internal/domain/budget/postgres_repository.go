package budget

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/db"
)

type PostgresRepository struct {
	queries *db.Queries
}

func NewPostgresRepository(dbConn *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: db.New(dbConn),
	}
}

func (r *PostgresRepository) Save(ctx context.Context, b Budget) error {
	dataJSON, err := json.Marshal(b.Data)
	if err != nil {
		return fmt.Errorf("marshal budget data: %w", err)
	}

	return r.queries.CreateBudget(ctx, db.CreateBudgetParams{
		ID:   b.ID,
		Name: b.Name,
		Data: dataJSON,
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBudget(ctx, id)
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Budget, error) {
	rows, err := r.queries.GetBudgets(ctx)
	if err != nil {
		return nil, fmt.Errorf("get budgets: %w", err)
	}

	budgets := make([]Budget, 0, len(rows))
	for _, row := range rows {
		var data any
		if err := json.Unmarshal(row.Data, &data); err != nil {
			return nil, fmt.Errorf("unmarshal budget data: %w", err)
		}

		budgets = append(budgets, Budget{
			ID:        row.ID,
			Name:      row.Name,
			Data:      data,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}

	return budgets, nil
}
