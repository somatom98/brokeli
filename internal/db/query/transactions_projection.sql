-- name: CreateTransaction :exec
INSERT INTO transactions (
    id, account_id, transaction_type, amount, currency, category, description, happened_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: ListTransactionsByAccount :many
WITH distributions AS (
    SELECT
        id,
        SUM(CASE WHEN origin = 'movement' AND user_id = 'system' THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY value_date ASC, id ASC) as system_amount,
        SUM(CASE WHEN origin = 'movement' AND user_id != 'system' THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY value_date ASC, id ASC) as other_amount
    FROM balance_updates
)
SELECT
    t.*,
    COALESCE(CASE WHEN (d.system_amount + d.other_amount) != 0 THEN (d.system_amount / (d.system_amount + d.other_amount))::TEXT ELSE '0' END, '0') as system_total_rate
FROM transactions t
LEFT JOIN distributions d ON t.id = d.id
WHERE t.account_id = $1
ORDER BY t.happened_at DESC;

-- name: ListTransactions :many
WITH distributions AS (
    SELECT
        id,
        SUM(CASE WHEN origin = 'movement' AND user_id = 'system' THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY value_date ASC, id ASC) as system_amount,
        SUM(CASE WHEN origin = 'movement' AND user_id != 'system' THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY value_date ASC, id ASC) as other_amount
    FROM balance_updates
)
SELECT
    t.*,
    COALESCE(CASE WHEN (d.system_amount + d.other_amount) != 0 THEN (d.system_amount / (d.system_amount + d.other_amount))::TEXT ELSE '0' END, '0') as system_total_rate
FROM transactions t
LEFT JOIN distributions d ON t.id = d.id
ORDER BY t.happened_at DESC;
