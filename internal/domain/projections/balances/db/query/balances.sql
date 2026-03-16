-- name: InsertBalanceUpdate :exec
INSERT INTO balances_projection (id, account_id, currency, amount, value_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO NOTHING;
