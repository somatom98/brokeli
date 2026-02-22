package transaction

import (
	"context"
	"fmt"

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
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.RegisterExpense(
		accountID,
		currency,
		amount,
		category,
		description,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
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
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.RegisterIncome(
		accountID,
		currency,
		amount,
		category,
		description,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
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
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.RegisterTransfer(
		fromAccountID,
		fromCurrency,
		fromAmount,
		toAccountID,
		toCurrency,
		toAmount,
		category,
		description,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})
}

func (d *Dispatcher) SetExpectedReimbursement(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.SetExpectedReimbursement(
		accountID,
		currency,
		amount,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})
}

func (d *Dispatcher) RegisterReimbursement(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	from string,
	currency values.Currency,
	amount decimal.Decimal,
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.RegisterReimbursement(
		accountID,
		from,
		currency,
		amount,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})
}

func (d *Dispatcher) RegisterDeposit(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.RegisterDeposit(
		accountID,
		currency,
		amount,
		category,
		description,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})
}

func (d *Dispatcher) RegisterWithdrawal(
	ctx context.Context,
	id uuid.UUID,
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
) error {
	aggr, err := d.es.GetAggregate(ctx, id)
	if err != nil {
		return fmt.Errorf("aggregate fetch failed: %w", err)
	}

	event, err := aggr.RegisterWithdrawal(
		accountID,
		currency,
		amount,
		category,
		description,
	)
	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}

	return d.es.Append(ctx, event_store.Record{
		AggregateID: aggr.ID,
		Version:     Version,
		Event:       event,
	})
}
