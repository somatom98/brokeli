CREATE TABLE balance_updates (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    currency TEXT NOT NULL,
    amount DECIMAL NOT NULL,
    user_id TEXT NOT NULL,
    origin TEXT NOT NULL,
    value_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
