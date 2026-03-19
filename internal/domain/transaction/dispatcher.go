package transaction

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type Dispatcher struct {
	es event_store.Store[*Transaction]
}

func NewDispatcher(
	es event_store.Store[*Transaction],
) *Dispatcher {
	return &Dispatcher{
		es: es,
	}
}

func (d *Dispatcher) RegisterExpense(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Transaction, version uint64) (event_store.Event, error) {
		return aggr.RegisterExpense(accountID, currency, amount, category, description, happenedAt)
	})
}

func (d *Dispatcher) RegisterIncome(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Transaction, version uint64) (event_store.Event, error) {
		return aggr.RegisterIncome(accountID, currency, amount, category, description, happenedAt)
	})
}

func (d *Dispatcher) RegisterTransfer(
	ctx context.Context,
	id uuid.UUID,
	fromAccountID uuid.UUID,
	fromCurrency values.Currency,
	fromAmount decimal.Decimal,
	toAccountID uuid.UUID,
	toCurrency values.Currency,
	toAmount decimal.Decimal,
	category string,
	description string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Transaction, version uint64) (event_store.Event, error) {
		return aggr.RegisterTransfer(fromAccountID, fromCurrency, fromAmount, toAccountID, toCurrency, toAmount, category, description, happenedAt)
	})
}

func (d *Dispatcher) SetExpectedReimbursement(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Transaction, version uint64) (event_store.Event, error) {
		return aggr.SetExpectedReimbursement(accountID, currency, amount, happenedAt)
	})
}

func (d *Dispatcher) RegisterReimbursement(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	from string,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Transaction, version uint64) (event_store.Event, error) {
		return aggr.RegisterReimbursement(accountID, from, currency, amount, category, description, happenedAt)
	})
}
