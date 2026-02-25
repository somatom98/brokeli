CREATE TABLE accounts_projection (
    id UUID PRIMARY KEY,
    balance JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP,
    closed_at TIMESTAMP
);
