package accounts

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type InMemoryRepository struct {
	accounts map[uuid.UUID]Account
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		accounts: make(map[uuid.UUID]Account),
	}
}

func (r *InMemoryRepository) CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error {
	acc := r.getOrCreate(id)
	if acc.CreatedAt == nil {
		acc.CreatedAt = &createdAt
	}
	r.accounts[id] = acc
	return nil
}

func (r *InMemoryRepository) CloseAccount(ctx context.Context, id uuid.UUID, closedAt time.Time) error {
	acc := r.getOrCreate(id)
	acc.ClosedAt = &closedAt
	r.accounts[id] = acc
	return nil
}

func (r *InMemoryRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error {
	acc := r.getOrCreate(id)

	if _, ok := acc.Balance[currency]; !ok {
		acc.Balance[currency] = decimal.Zero
	}

	acc.Balance[currency] = acc.Balance[currency].Add(amount)
	r.accounts[id] = acc

	return nil
}

func (r *InMemoryRepository) SetExpectedReimbursement(ctx context.Context, id uuid.UUID, amount decimal.Decimal, currency values.Currency) error {
	acc := r.getOrCreate(id)

	acc.ExpectedReimbursements[currency] = amount
	r.accounts[id] = acc

	return nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) (map[uuid.UUID]Account, error) {
	return r.accounts, nil
}

func (r *InMemoryRepository) getOrCreate(id uuid.UUID) Account {
	if _, ok := r.accounts[id]; !ok {
		r.accounts[id] = NewAccount()
	}

	return r.accounts[id]
}
