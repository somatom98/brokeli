package event_store

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemoryStore[A Aggregate] struct {
	mu       sync.RWMutex
	locks    map[uuid.UUID]*sync.Mutex
	events   map[uuid.UUID][]Record
	handlers []SubscribeHandler
	new      func(uuid.UUID) A
}

func NewInMemory[A Aggregate](new func(uuid.UUID) A) *InMemoryStore[A] {
	return &InMemoryStore[A]{
		locks:    make(map[uuid.UUID]*sync.Mutex),
		events:   make(map[uuid.UUID][]Record),
		handlers: make([]SubscribeHandler, 0),
		new:      new,
	}
}

func (s *InMemoryStore[A]) getLock(id uuid.UUID) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	if lock, ok := s.locks[id]; ok {
		return lock
	}

	lock := &sync.Mutex{}
	s.locks[id] = lock
	return lock
}

func (s *InMemoryStore[A]) Subscribe(ctx context.Context, handler SubscribeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers = append(s.handlers, handler)
}

func (s *InMemoryStore[A]) Execute(ctx context.Context, id uuid.UUID, fn func(aggr A, version uint64) (Event, error)) error {
	lock := s.getLock(id)
	lock.Lock()
	defer lock.Unlock()

	aggr, version, err := s.GetAggregate(ctx, id)
	if err != nil {
		return err
	}

	e, err := fn(aggr, version)
	if err != nil {
		return err
	}

	if e == nil {
		return nil
	}

	return s.Append(ctx, Record{
		AggregateID: id,
		Version:     version + 1,
		Event: simpleEvent{
			eventType: e.Type(),
			content:   e.Content(),
		},
	})
}

type simpleEvent struct {
	eventType string
	content   any
}

func (e simpleEvent) Type() string {
	return e.eventType
}

func (e simpleEvent) Content() any {
	return e.content
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
