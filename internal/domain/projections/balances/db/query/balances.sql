-- name: InsertBalanceUpdate :exec
INSERT INTO balances_projection (id, account_id, currency, amount, value_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO NOTHING;

-- name: GetBalancesByAccount :many
SELECT DATE_TRUNC('month', value_date)::TIMESTAMP AS month, currency, SUM(amount)::TEXT AS amount
FROM balances_projection
WHERE account_id = $1
GROUP BY month, currency
ORDER BY month DESC;

-- name: GetAllBalances :many
SELECT DATE_TRUNC('month', value_date)::TIMESTAMP AS month, currency, SUM(amount)::TEXT AS amount
FROM balances_projection
GROUP BY month, currency
ORDER BY month DESC;
