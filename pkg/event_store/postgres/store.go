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
	db           *sql.DB
	new          func(uuid.UUID) A
	eventFactory map[string]func() interface{}
	subscribers  []chan event_store.Record
	mu           sync.RWMutex
}

var _ event_store.Store[event_store.Aggregate] = &PostgresStore[event_store.Aggregate]{}

func NewPostgresStore[A event_store.Aggregate](
	db *sql.DB,
	new func(uuid.UUID) A,
	eventFactory map[string]func() interface{},
) (*PostgresStore[A], error) {
	if _, err := db.Exec(Schema); err != nil {
		return nil, fmt.Errorf("failed to ensure schema: %w", err)
	}

	return &PostgresStore[A]{
		db:           db,
		new:          new,
		eventFactory: eventFactory,
		subscribers:  make([]chan event_store.Record, 0),
	}, nil
}

func (s *PostgresStore[A]) Subscribe(ctx context.Context) <-chan event_store.Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan event_store.Record, 100)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *PostgresStore[A]) Append(ctx context.Context, record event_store.Record) error {
	aggregateType := s.getAggregateType()

	// 1. Serialize Event
	eventData, err := json.Marshal(record.Event.Content())
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// 2. Serialize Snapshot (The aggregate state AFTER the event)
	// We need the aggregate to serialize it.
	// But `record` only has the event.
	// We need to apply the event to the current aggregate state to get the new state.
	// OR, the caller should have passed the *Aggregate*?
	// The `Store` interface `Append` only takes `Record`.
	// `Record` has `AggregateID`.
	
	// CRITICAL ISSUE: The `Append` signature only takes `Record`.
	// To update the snapshot, we need the *Resulting Aggregate*.
	// We have two options:
	// A. Load current snapshot, apply event, save new snapshot. (Read-Modify-Write)
	// B. Require `Append` to take the Aggregate (Change Interface).
	// C. Assume the `Record` comes from a Dispatcher that *just* modified the Aggregate,
	//    so maybe we don't have the aggregate here?
	
	// Actually, in `Dispatcher`, we do:
	// aggr.Create(...) -> returns event
	// es.Append(..., record)
	// The `aggr` variable in Dispatcher HAS the new state.
	// But `Append` doesn't receive `aggr`.
	
	// Option A (Read-Modify-Write) inside Append is safest but slower (extra read).
	// But wait, `Append` is usually called *after* we've validated logic on the aggregate.
	// If we reload the aggregate from DB, apply event, we are duplicating work?
	// No, we can just load the *latest snapshot* (fast), apply *this one event* (fast), and save.
	
	// Let's do Option A for now as it fits the interface.
	
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert Event
	eventID := uuid.New()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO events (id, aggregate_id, aggregate_type, version, event_type, event_data)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, eventID, record.AggregateID, aggregateType, record.Version, record.Event.Type(), eventData)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	// Re-calculate Aggregate State for Snapshot
	// We need to fetch the *previous* state to apply the *new* event.
	// Optimally, we would have the aggregate passed in.
	// Since we don't, we fetch the snapshot.
	
	// Fetch current snapshot
	var stateJSON []byte
	err = tx.QueryRowContext(ctx, `SELECT state FROM snapshots WHERE aggregate_id = $1`, record.AggregateID).Scan(&stateJSON)
	
	aggregate := s.new(record.AggregateID)
	if err == nil {
		// Snapshot exists
		if err := json.Unmarshal(stateJSON, aggregate); err != nil {
			return fmt.Errorf("failed to unmarshal snapshot state: %w", err)
		}
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed to fetch snapshot: %w", err)
	}
	// If ErrNoRows, we start with new(id) which is correct.

	// Apply the NEW event to the aggregate
	// We need to put the record in a slice
	err = aggregate.Hydrate([]event_store.Record{record})
	if err != nil {
		return fmt.Errorf("failed to hydrate aggregate for snapshot: %w", err)
	}

	// Serialize new state
	newStateJSON, err := json.Marshal(aggregate)
	if err != nil {
		return fmt.Errorf("failed to marshal new snapshot state: %w", err)
	}

	// Upsert Snapshot
	_, err = tx.ExecContext(ctx, `
		INSERT INTO snapshots (aggregate_id, aggregate_type, version, state, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (aggregate_id) DO UPDATE
		SET version = $3, state = $4, updated_at = NOW()
	`, record.AggregateID, aggregateType, record.Version, newStateJSON)
	if err != nil {
		return fmt.Errorf("failed to upsert snapshot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Notify subscribers (best effort)
	s.publishToSubscribers(record)

	return nil
}

func (s *PostgresStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var zero A
	
	// Try to get from Snapshot
	var stateJSON []byte
	err := s.db.QueryRowContext(ctx, `SELECT state FROM snapshots WHERE aggregate_id = $1`, id).Scan(&stateJSON)
	
	if err == sql.ErrNoRows {
		// No snapshot? Return new empty aggregate.
		// NOTE: In a pure event sourcing world, we might want to check if events exist.
		// But here we enforce snapshots. If no snapshot, it's a new aggregate.
		return s.new(id), nil
	}
	if err != nil {
		return zero, fmt.Errorf("failed to fetch snapshot: %w", err)
	}

	aggregate := s.new(id)
	if err := json.Unmarshal(stateJSON, aggregate); err != nil {
		return zero, fmt.Errorf("failed to unmarshal snapshot: %w", err)
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
