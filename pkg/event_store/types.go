package event_store

import (
	"context"

	"github.com/gofrs/uuid"
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

type Store[A Aggregate] interface {
	Subscribe(ctx context.Context) <-chan Record
	GetAggregate(ctx context.Context, id uuid.UUID) (A, error)
	Append(ctx context.Context, record Record) error
}
