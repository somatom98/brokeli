package event_store

import (
	"context"

	"github.com/google/uuid"
)

type Event interface {
	Type() string
	Content() any
}

type Record struct {
	AggregateID uuid.UUID
	Version     uint64
	Event
}

type Aggregate interface {
	Hydrate(records []Record) error
}

type SubscribeHandler func(ctx context.Context, record Record) error

type Store[A Aggregate] interface {
	Subscribe(ctx context.Context, handler SubscribeHandler)
	GetAggregate(ctx context.Context, id uuid.UUID) (A, uint64, error)
	Append(ctx context.Context, record Record) error
}
