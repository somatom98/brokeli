package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/pkg/event_store"
)

// MockAggregate for testing
type MockAggregate struct {
	ID    uuid.UUID
	State string
}

func (a *MockAggregate) Hydrate(records []event_store.Record) error {
	for _, r := range records {
		var content string
		data, _ := json.Marshal(r.Event.Content())
		json.Unmarshal(data, &content)
		a.State = content
	}
	return nil
}

type MockEvent struct {
	typ     string
	content string
}

func (e MockEvent) Type() string { return e.typ }
func (e MockEvent) Content() any { return e.content }

func TestPostgresStore_Concurrency_UniqueVersion(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("Skipping integration test: DB_DSN not set")
	}

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Ensure clean state for this test
	// Ideally we run migration here
	_, err = db.Exec(Schema)
	require.NoError(t, err)

	id := uuid.New()
	store, err := NewPostgresStore[*MockAggregate](
		db,
		func(uid uuid.UUID) *MockAggregate { return &MockAggregate{ID: uid} },
		nil,
	)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()
	version := uint64(1)

	rec1 := event_store.Record{
		AggregateID: id,
		Version:     version,
		Event:       MockEvent{typ: "TEST", content: "A"},
	}
	rec2 := event_store.Record{
		AggregateID: id,
		Version:     version,
		Event:       MockEvent{typ: "TEST", content: "B"},
	}

	// First append should succeed
	err = store.Append(ctx, rec1)
	require.NoError(t, err)

	// Second append with SAME version should FAIL
	err = store.Append(ctx, rec2)
	assert.Error(t, err, "Expected error when appending event with duplicate version")
}
