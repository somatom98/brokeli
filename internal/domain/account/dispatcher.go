package account

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type Dispatcher struct {
	es event_store.Store[*Account]
}

func NewDispatcher(
	es event_store.Store[*Account],
) *Dispatcher {
	return &Dispatcher{
		es: es,
	}
}

func (d *Dispatcher) Open(
	ctx context.Context,
	id uuid.UUID,
	name string,
	currency values.Currency,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Account, version uint64) (event_store.Event, error) {
		return aggr.Open(name, currency, happenedAt)
	})
}

func (d *Dispatcher) UpdateName(
	ctx context.Context,
	id uuid.UUID,
	name string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Account, version uint64) (event_store.Event, error) {
		return aggr.UpdateName(name, happenedAt)
	})
}

func (d *Dispatcher) Deposit(
	ctx context.Context,
	id uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	user string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Account, version uint64) (event_store.Event, error) {
		return aggr.Deposit(currency, amount, category, description, user, happenedAt)
	})
}

func (d *Dispatcher) Withdraw(
	ctx context.Context,
	id uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	user string,
	happenedAt time.Time,
) error {
	return d.es.Execute(ctx, id, func(aggr *Account, version uint64) (event_store.Event, error) {
		return aggr.Withdraw(currency, amount, category, description, user, happenedAt)
	})
}
