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
	"github.com/somatom98/brokeli/pkg/event_store/postgres/db"
)

type PostgresStore[A event_store.Aggregate] struct {
	db            *sql.DB
	queries       *db.Queries
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
	dbConn *sql.DB,
	new func(uuid.UUID) A,
	eventFactory map[string]func() any,
) (*PostgresStore[A], error) {
	store := &PostgresStore[A]{
		db:           dbConn,
		queries:      db.New(dbConn),
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

	err = s.queries.AppendEvent(ctx, db.AppendEventParams{
		ID:            uuid.New(),
		AggregateID:   record.AggregateID,
		AggregateType: aggregateType,
		Version:       int64(record.Version),
		EventType:     record.Type(),
		EventData:     eventData,
	})
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	s.publishToSubscribers(record)

	return nil
}

func (s *PostgresStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var zero A

	events, err := s.queries.GetEvents(ctx, id)
	if err != nil {
		return zero, fmt.Errorf("failed to query events: %w", err)
	}

	var records []event_store.Record
	for _, row := range events {
		factory, ok := s.eventFactory[row.EventType]
		if !ok {
			return zero, fmt.Errorf("unknown event type: %s", row.EventType)
		}

		eventPtr := factory()
		if err := json.Unmarshal(row.EventData, eventPtr); err != nil {
			return zero, fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		// Dereference if it's a pointer to get the value
		content := reflect.ValueOf(eventPtr).Elem().Interface()

		records = append(records, event_store.Record{
			AggregateID: id,
			Version:     uint64(row.Version),
			Event: event{
				EventType:    row.EventType,
				EventContent: content,
			},
		})
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
		ch <- record
	}
}
