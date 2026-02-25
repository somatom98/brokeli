package postgres

const Schema = `
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    event_data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_aggregate_id ON events (aggregate_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_events_aggregate_id_version ON events (aggregate_id, version);
`
