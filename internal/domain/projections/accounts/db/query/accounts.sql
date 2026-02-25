-- name: CreateAccount :exec
INSERT INTO accounts_projection (id, created_at, balance)
VALUES ($1, $2, '{}')
ON CONFLICT (id) DO UPDATE SET created_at = EXCLUDED.created_at;

-- name: CloseAccount :exec
UPDATE accounts_projection
SET closed_at = $2
WHERE id = $1;

-- name: GetAccountBalanceForUpdate :one
SELECT balance FROM accounts_projection WHERE id = $1 FOR UPDATE;

-- name: UpdateAccountBalance :exec
UPDATE accounts_projection SET balance = $2 WHERE id = $1;

-- name: GetAllAccounts :many
SELECT id, balance, created_at, closed_at FROM accounts_projection;

-- name: UpsertPlaceholderAccount :exec
INSERT INTO accounts_projection (id, balance) VALUES ($1, $2)
ON CONFLICT (id) DO NOTHING;
