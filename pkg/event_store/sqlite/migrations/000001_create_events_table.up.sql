CREATE TABLE events (
    id TEXT PRIMARY KEY,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    version INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    event_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(aggregate_type, aggregate_id, version)
);

CREATE INDEX idx_aggregate_type_id ON events(aggregate_type, aggregate_id);
CREATE INDEX idx_created_at ON events(created_at);