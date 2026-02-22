package account

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type Account struct {
	ID        uuid.UUID
	Balances  map[values.Currency]decimal.Decimal
	CreatedAt time.Time
	ClosedAt  *time.Time
}

func New(id uuid.UUID) *Account {
	return &Account{
		ID:       id,
		Balances: make(map[values.Currency]decimal.Decimal),
	}
}
