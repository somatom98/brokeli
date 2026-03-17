-- name: CreateBudget :exec
INSERT INTO budgets (id, name, data, created_at, updated_at)
VALUES ($1, $2, $3, $4, $4)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    data = EXCLUDED.data,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE id = $1;

-- name: GetBudgets :many
SELECT * FROM budgets
ORDER BY name ASC;

-- name: GetBudgetByID :one
SELECT * FROM budgets
WHERE id = $1;
