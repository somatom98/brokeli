package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type PostgresStore[A event_store.Aggregate] struct {
	db            *sql.DB
	new           func(uuid.UUID) A
	eventFactory  map[string]func() any
	subscribers   []chan event_store.Record
	mu            sync.RWMutex
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
	db *sql.DB,
	new func(uuid.UUID) A,
	eventFactory map[string]func() any,
) (*PostgresStore[A], error) {
	if _, err := db.Exec(Schema); err != nil {
		return nil, fmt.Errorf("failed to ensure schema: %w", err)
	}

	store := &PostgresStore[A]{
		db:           db,
		new:          new,
		eventFactory: eventFactory,
		subscribers:  make([]chan event_store.Record, 0),
	}
	store.aggregateType = store.getAggregateType()
	return store, nil
}

func (s *PostgresStore[A]) Subscribe(ctx context.Context) <-chan event_store.Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan event_store.Record, 100)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *PostgresStore[A]) Append(ctx context.Context, record event_store.Record) error {
	aggregateType := s.aggregateType

	eventData, err := json.Marshal(record.Content())
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	eventID := uuid.New()
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO events (id, aggregate_id, aggregate_type, version, event_type, event_data)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, eventID, record.AggregateID, aggregateType, record.Version, record.Type(), eventData)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	s.publishToSubscribers(record)

	return nil
}

func (s *PostgresStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var zero A

	rows, err := s.db.QueryContext(ctx, `
		SELECT version, event_type, event_data
		FROM events
		WHERE aggregate_id = $1
		ORDER BY version ASC
	`, id)
	if err != nil {
		return zero, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var records []event_store.Record
	for rows.Next() {
		var version uint64
		var eventType string
		var eventData []byte

		if err := rows.Scan(&version, &eventType, &eventData); err != nil {
			return zero, fmt.Errorf("failed to scan event: %w", err)
		}

		factory, ok := s.eventFactory[eventType]
		if !ok {
			return zero, fmt.Errorf("unknown event type: %s", eventType)
		}

		eventPtr := factory()
		if err := json.Unmarshal(eventData, eventPtr); err != nil {
			return zero, fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		// Dereference if it's a pointer to get the value
		content := reflect.ValueOf(eventPtr).Elem().Interface()

		records = append(records, event_store.Record{
			AggregateID: id,
			Version:     version,
			Event: event{
				EventType:    eventType,
				EventContent: content,
			},
		})
	}

	if err = rows.Err(); err != nil {
		return zero, fmt.Errorf("rows error: %w", err)
	}

	aggregate := s.new(id)
	if len(records) > 0 {
		if err := aggregate.Hydrate(records); err != nil {
			return zero, fmt.Errorf("failed to hydrate aggregate: %w", err)
		}
	}

	return aggregate, nil
}

func (s *PostgresStore[A]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ch := range s.subscribers {
		close(ch)
	}
	s.subscribers = nil
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

func (s *PostgresStore[A]) publishToSubscribers(record event_store.Record) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.subscribers {
		select {
		case ch <- record:
		default:
			// Drop if full
		}
	}
}
