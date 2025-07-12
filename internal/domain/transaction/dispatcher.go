package transaction

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
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

func (d *Dispatcher) CreateExpense(ctx context.Context, id uuid.UUID, cmd CreateExpense) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.HandleCreateExpense(cmd)
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

func (d *Dispatcher) CreateIncome(ctx context.Context, id uuid.UUID, cmd CreateIncome) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.HandleCreateIncome(cmd)
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

func (d *Dispatcher) CreateTransfer(ctx context.Context, id uuid.UUID, cmd CreateTransfer) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.HandleCreateTransfer(cmd)
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
