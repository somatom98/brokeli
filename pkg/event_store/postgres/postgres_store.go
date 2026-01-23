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
	"github.com/somatom98/brokeli/pkg/event_store/postgres/generated"
)

type PostgresStore[A event_store.Aggregate] struct {
	db          *sql.DB
	queries     *generated.Queries
	new         func(uuid.UUID) A
	subscribers []chan event_store.Record
	mu          sync.RWMutex
	stopCh      chan struct{}
	stopped     bool
}

var _ event_store.Store[event_store.Aggregate] = &PostgresStore[event_store.Aggregate]{}

func NewPostgresStore[A event_store.Aggregate](db *sql.DB, new func(uuid.UUID) A) *PostgresStore[A] {
	return &PostgresStore[A]{
		db:          db,
		queries:     generated.New(db),
		new:         new,
		subscribers: make([]chan event_store.Record, 0),
		stopCh:      make(chan struct{}),
	}
}

func (s *PostgresStore[A]) Subscribe(ctx context.Context) <-chan event_store.Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		ch := make(chan event_store.Record)
		close(ch)
		return ch
	}

	ch := make(chan event_store.Record, 100)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *PostgresStore[A]) Append(ctx context.Context, record event_store.Record) error {
	aggregateType := s.getAggregateType()

	eventData, err := json.Marshal(record.Event.Content())
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	eventID := uuid.New()

	_, err = s.queries.InsertEvent(ctx, generated.InsertEventParams{
		ID:            eventID.String(),
		AggregateType: aggregateType,
		AggregateID:   record.AggregateID.String(),
		Version:       int64(record.Version),
		EventType:     record.Event.Type(),
		EventData:     string(eventData),
	})

	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	s.publishToSubscribers(record)

	return nil
}

func (s *PostgresStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var zero A
	aggregateType := s.getAggregateType()

	events, err := s.queries.GetEventsByAggregateID(ctx, generated.GetEventsByAggregateIDParams{
		AggregateType: aggregateType,
		AggregateID:   id.String(),
	})
	if err != nil {
		return zero, fmt.Errorf("failed to get events: %w", err)
	}

	var records []event_store.Record
	for _, event := range events {
		var content map[string]interface{}
		if err := json.Unmarshal([]byte(event.EventData), &content); err != nil {
			return zero, fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		genericEvent := &genericEvent{
			eventType: event.EventType,
			content:   content,
		}

		records = append(records, event_store.Record{
			AggregateID: id,
			Version:     uint64(event.Version),
			Event:       genericEvent,
		})
	}

	aggregate := s.new(id)
	if err := aggregate.Hydrate(records); err != nil {
		return zero, fmt.Errorf("failed to hydrate aggregate: %w", err)
	}

	return aggregate, nil
}

func (s *PostgresStore[A]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return nil
	}

	s.stopped = true
	close(s.stopCh)

	for _, ch := range s.subscribers {
		close(ch)
	}
	s.subscribers = nil

	return s.db.Close()
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
		}
	}
}

type genericEvent struct {
	eventType string
	content   interface{}
}

func (e *genericEvent) Type() string {
	return e.eventType
}

func (e *genericEvent) Content() any {
	return e.content
}