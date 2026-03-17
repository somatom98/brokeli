-- name: InsertBalanceUpdate :exec
INSERT INTO balance_updates (id, account_id, currency, amount, user_id, origin, value_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO NOTHING;

-- name: GetBalancesByAccount :many
SELECT DATE_TRUNC('month', value_date)::TIMESTAMP AS month, currency, SUM(amount)::TEXT AS amount
FROM balance_updates
WHERE account_id = $1
GROUP BY month, currency
ORDER BY month DESC;

-- name: GetAllBalances :many
SELECT DATE_TRUNC('month', value_date)::TIMESTAMP AS month, currency, SUM(amount)::TEXT AS amount
FROM balance_updates
GROUP BY month, currency
ORDER BY month DESC;

-- name: GetAccountDistributions :many
SELECT 
    id, account_id, currency, amount, user_id, value_date,
    SUM(CASE WHEN user_id = 'system' THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY value_date ASC, id ASC)::TEXT as system_amount,
    SUM(CASE WHEN user_id != 'system' THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY value_date ASC, id ASC)::TEXT as other_amount
FROM balance_updates
WHERE account_id = $1 AND origin = 'movement'
ORDER BY value_date DESC, id DESC;
