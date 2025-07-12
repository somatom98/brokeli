package event_store

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

var ErrEventsNotFound = errors.New("no_events_found")

type InMemoryStore[A Aggregate] struct {
	records map[uuid.UUID][]Record
	new     func(uuid.UUID) A
	ch      chan Record
}

var _ Store[Aggregate] = &InMemoryStore[Aggregate]{}

func NewInMemory[A Aggregate](new func(uuid.UUID) A) *InMemoryStore[A] {
	return &InMemoryStore[A]{
		records: make(map[uuid.UUID][]Record),
		new:     new,
		ch:      make(chan Record),
	}
}

func (s *InMemoryStore[A]) Subscribe(ctx context.Context) <-chan Record {
	return s.ch
}

func (s *InMemoryStore[A]) Append(ctx context.Context, record Record) error {
	if _, ok := s.records[record.AggregateID]; !ok {
		s.records[record.AggregateID] = []Record{}
	}
	s.records[record.AggregateID] = append(s.records[record.AggregateID], record)
	s.ch <- record
	return nil
}

func (s *InMemoryStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var aggregate A

	records, ok := s.records[id]
	if !ok {
		aggregate = s.new(id)
	}

	err := aggregate.Hydrate(records)
	if err != nil {
		return aggregate, fmt.Errorf("hydration failed: %w", err)
	}

	return aggregate, nil
}
