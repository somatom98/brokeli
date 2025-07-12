package accounts

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type InMemoryRepository struct {
	accounts map[uuid.UUID]Account
}

func (r *InMemoryRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error {
	if _, ok := r.accounts[id]; !ok {
		r.accounts[id] = NewAccount()
	}

	if _, ok := r.accounts[id].Balance[currency]; !ok {
		r.accounts[id].Balance[currency] = decimal.Zero
	}

	r.accounts[id].Balance[currency].Add(amount)

	return nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]Account, error) {
	accounts := []Account{}
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}
