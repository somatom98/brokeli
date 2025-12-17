package account

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type Dispatcher struct {
	es event_store.Store[*Account]
}

func NewDispatcher(es event_store.Store[*Account]) *Dispatcher {
	return &Dispatcher{
		es: es,
	}
}

func (d *Dispatcher) CreateAccount(ctx context.Context, id uuid.UUID, createdAt time.Time) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.Create(createdAt)
	if err != nil {
		return err
	}

	d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})

	return nil
}

func (d *Dispatcher) DepositMoney(ctx context.Context, id uuid.UUID, user string, currency values.Currency, amount decimal.Decimal, time time.Time) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.Deposit(user, currency, amount, time)
	if err != nil {
		return err
	}

	d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})

	return nil
}

func (d *Dispatcher) CloseAccount(ctx context.Context, id uuid.UUID, time time.Time) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.Close(time)
	if err != nil {
		return err
	}

	d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})

	return nil
}
