package event_store

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

var ErrEventsNotFound = errors.New("no_events_found")

type InMemoryStore[A Aggregate] struct {
	Records map[uuid.UUID][]Record
	New     func(uuid.UUID) A
}

var _ Store[Aggregate] = &InMemoryStore[Aggregate]{}

func NewInMemory[A Aggregate](new func(uuid.UUID) A) *InMemoryStore[A] {
	return &InMemoryStore[A]{
		Records: make(map[uuid.UUID][]Record),
		New:     new,
	}
}

func (s *InMemoryStore[A]) Append(ctx context.Context, id uuid.UUID, record Record) error {
	if _, ok := s.Records[id]; !ok {
		s.Records[id] = []Record{}
	}
	s.Records[id] = append(s.Records[id], record)
	return nil
}

func (s *InMemoryStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var aggregate A

	records, ok := s.Records[id]
	if !ok {
		aggregate = s.New(id)
	}

	err := aggregate.Hydrate(records)
	if err != nil {
		return aggregate, fmt.Errorf("hydration failed: %w", err)
	}

	return aggregate, nil
}
