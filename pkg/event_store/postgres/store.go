package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/pkg/event_store"
	"github.com/somatom98/brokeli/pkg/event_store/postgres/db"
)

type PostgresStore[A event_store.Aggregate] struct {
	db            *sql.DB
	queries       *db.Queries
	new           func(uuid.UUID) A
	eventFactory  map[string]func() any
	mu            sync.RWMutex
	handlers      []event_store.SubscribeHandler
	aggregateType string
}

type event struct {
	EventType    string
	EventContent any
}

func (e event) Type() string {
	return e.EventType
}

func (e event) Content() any {
	return e.EventContent
}

var _ event_store.Store[event_store.Aggregate] = &PostgresStore[event_store.Aggregate]{}

func NewPostgresStore[A event_store.Aggregate](
	dbConn *sql.DB,
	new func(uuid.UUID) A,
	eventFactory map[string]func() any,
) (*PostgresStore[A], error) {
	store := &PostgresStore[A]{
		db:           dbConn,
		queries:      db.New(dbConn),
		new:          new,
		eventFactory: eventFactory,
		handlers:     make([]event_store.SubscribeHandler, 0),
	}
	store.aggregateType = store.getAggregateType()
	return store, nil
}

func (s *PostgresStore[A]) RunRelay(ctx context.Context) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_ = s.handleEvents(ctx)
		}
	}
}

func (s *PostgresStore[A]) Subscribe(ctx context.Context, handler event_store.SubscribeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers = append(s.handlers, handler)
}

func (s *PostgresStore[A]) Append(ctx context.Context, record event_store.Record) error {
	aggregateType := s.aggregateType

	eventData, err := json.Marshal(record.Content())
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	params := db.AppendEventParams{
		ID:            uuid.New(),
		AggregateID:   record.AggregateID,
		AggregateType: aggregateType,
		Version:       int64(record.Version),
		EventType:     record.Type(),
		EventData:     eventData,
	}

	err = qtx.AppendEvent(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	err = qtx.AppendToOutbox(ctx, db.AppendToOutboxParams(params))
	if err != nil {
		return fmt.Errorf("failed to insert outbox event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PostgresStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, uint64, error) {
	var zero A

	events, err := s.queries.GetEvents(ctx, id)
	if err != nil {
		return zero, 0, fmt.Errorf("failed to query events: %w", err)
	}

	aggregate := s.new(id)
	var version uint64 = 0
	for _, row := range events {
		version++
		if version != uint64(row.Version) {
			return zero, version, fmt.Errorf("invalid version number for aggregate %s: %v, expected %v", id, row.Version, version)
		}

		factory, ok := s.eventFactory[row.EventType]
		if !ok {
			return zero, version, fmt.Errorf("unknown event type: %s", row.EventType)
		}

		eventPtr := factory()
		if err := json.Unmarshal(row.EventData, eventPtr); err != nil {
			return zero, version, fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		// Dereference if it's a pointer to get the value
		content := reflect.ValueOf(eventPtr).Elem().Interface()

		records := []event_store.Record{
			{
				AggregateID: id,
				Version:     uint64(row.Version),
				Event: event{
					EventType:    row.EventType,
					EventContent: content,
				},
			},
		}

		if err := aggregate.Hydrate(records); err != nil {
			return zero, version, fmt.Errorf("failed to hydrate aggregate: %w", err)
		}
	}

	return aggregate, version, nil
}

func (s *PostgresStore[A]) Close() error {
	return nil
}

func (s *PostgresStore[A]) handleEvents(ctx context.Context) error {
	rows, err := s.queries.GetOutboxEvents(ctx, 10)
	if err != nil {
		return fmt.Errorf("failed to get outbox events: %w", err)
	}

	if len(rows) == 0 {
		return nil
	}

	for _, row := range rows {
		factory, ok := s.eventFactory[row.EventType]
		if !ok {
			return fmt.Errorf("event factory not found for event %s", row.EventType)
		}

		eventPtr := factory()
		if err := json.Unmarshal(row.EventData, eventPtr); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		content := reflect.ValueOf(eventPtr).Elem().Interface()

		record := event_store.Record{
			AggregateID: row.AggregateID,
			Version:     uint64(row.Version),
			Event: event{
				EventType:    row.EventType,
				EventContent: content,
			},
		}

		s.mu.RLock()
		handlers := s.handlers
		s.mu.RUnlock()

		for _, h := range handlers {
			_ = h(ctx, record)
		}

		err = s.queries.DeleteOutboxEvent(ctx, row.ID)
		if err != nil {
			return fmt.Errorf("failed to delete outbox %s: %w", row.ID, err)
		}
	}
	return nil
}

func (s *PostgresStore[A]) getAggregateType() string {
	aggregate := s.new(uuid.New())
	aggregateType := reflect.TypeOf(aggregate)
	if aggregateType.Kind() == reflect.Ptr {
		aggregateType = aggregateType.Elem()
	}

	parts := strings.Split(aggregateType.String(), ".")
	return parts[len(parts)-1]
}
