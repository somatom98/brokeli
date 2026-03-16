CREATE TABLE balances_projection (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    currency TEXT NOT NULL,
    amount DECIMAL NOT NULL,
    value_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
