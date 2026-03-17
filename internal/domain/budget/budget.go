package budget

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Data      any       `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Save(ctx context.Context, b Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]Budget, error)
}
