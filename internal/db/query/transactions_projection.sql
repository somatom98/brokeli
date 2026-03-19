-- name: CreateTransaction :exec
INSERT INTO transactions (
    id, account_id, transaction_type, amount, currency, category, description, happened_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: ListTransactions :many
WITH distributions AS (
    SELECT
        id,
        SUM(CASE WHEN transaction_type IN ('TRANSFER') THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY happened_at ASC, id ASC) as system_amount,
        SUM(CASE WHEN transaction_type IN ('DEPOSIT', 'WITHDRAWAL', 'INCOME', 'EXPENSE') THEN amount ELSE 0 END) OVER (PARTITION BY account_id, currency ORDER BY happened_at ASC, id ASC) as other_amount
    FROM transactions
)
SELECT
    t.id, t.account_id, t.transaction_type, t.amount, t.currency, t.category, t.description, t.happened_at, t.created_at,
    COALESCE(CASE 
        WHEN d.system_amount + d.other_amount != 0 THEN ROUND(d.system_amount::DECIMAL / (d.system_amount + d.other_amount)::DECIMAL, 4)::TEXT ELSE '0' END, '0') as system_total_rate
FROM transactions t
JOIN distributions d ON t.id = d.id
WHERE
    (t.happened_at >= sqlc.narg('start_date') OR sqlc.narg('start_date') IS NULL) AND
    (t.happened_at <= sqlc.narg('end_date') OR sqlc.narg('end_date') IS NULL) AND
    (t.account_id = ANY(sqlc.narg('account_ids')::UUID[]) OR sqlc.narg('account_ids') IS NULL) AND
    (t.transaction_type = sqlc.narg('transaction_type') OR sqlc.narg('transaction_type') IS NULL)
ORDER BY t.happened_at DESC, t.id DESC;

-- name: ListCategories :many
SELECT DISTINCT category FROM transactions ORDER BY category;
