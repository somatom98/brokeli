-- name: AppendEvent :exec
INSERT INTO events (id, aggregate_id, aggregate_type, version, event_type, event_data)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetEvents :many
SELECT version, event_type, event_data
FROM events
WHERE aggregate_id = $1
ORDER BY version ASC;
