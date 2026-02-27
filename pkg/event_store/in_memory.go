package event_store

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemoryStore[A Aggregate] struct {
	mu        sync.RWMutex
	events    map[uuid.UUID][]Record
	handlers  []SubscribeHandler
	new       func(uuid.UUID) A
}

func NewInMemory[A Aggregate](new func(uuid.UUID) A) *InMemoryStore[A] {
	return &InMemoryStore[A]{
		events:   make(map[uuid.UUID][]Record),
		handlers: make([]SubscribeHandler, 0),
		new:      new,
	}
}

func (s *InMemoryStore[A]) Subscribe(ctx context.Context, handler SubscribeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers = append(s.handlers, handler)
}

func (s *InMemoryStore[A]) Append(ctx context.Context, record Record) error {
	s.mu.Lock()
	s.events[record.AggregateID] = append(s.events[record.AggregateID], record)
	handlers := make([]SubscribeHandler, len(s.handlers))
	copy(handlers, s.handlers)
	s.mu.Unlock()

	for _, h := range handlers {
		_ = h(ctx, record)
	}

	return nil
}

func (s *InMemoryStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	aggregate := s.new(id)
	records, ok := s.events[id]
	if !ok {
		return aggregate, 0, nil
	}

	if err := aggregate.Hydrate(records); err != nil {
		var zero A
		return zero, 0, err
	}

	return aggregate, uint64(len(records)), nil
}
