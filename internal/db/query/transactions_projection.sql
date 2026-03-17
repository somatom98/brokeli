-- name: CreateTransaction :exec
INSERT INTO transactions (
    id, account_id, transaction_type, amount, currency, category, description, happened_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: ListTransactionsByAccount :many
SELECT * 
FROM transactions
WHERE account_id = $1
ORDER BY happened_at DESC;

-- name: ListTransactions :many
SELECT * 
FROM transactions
ORDER BY happened_at DESC;
