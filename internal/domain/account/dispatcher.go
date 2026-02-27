package account

import (
	"context"
	"fmt"

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
) error {
	aggr, version, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.Open(
		name,
		currency,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     version + 1,
		Event:       event,
	})
}

func (d *Dispatcher) UpdateName(
	ctx context.Context,
	id uuid.UUID,
	name string,
) error {
	aggr, version, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.UpdateName(
		name,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     version + 1,
		Event:       event,
	})
}

func (d *Dispatcher) Deposit(
	ctx context.Context,
	id uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	user string,
) error {
	aggr, version, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.Deposit(
		currency,
		amount,
		user,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     version + 1,
		Event:       event,
	})
}

func (d *Dispatcher) Withdraw(
	ctx context.Context,
	id uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	user string,
) error {
	aggr, version, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.Withdraw(
		currency,
		amount,
		user,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     version + 1,
		Event:       event,
	})
}
