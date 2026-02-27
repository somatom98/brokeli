-- name: AppendEvent :exec
INSERT INTO events (id, aggregate_id, aggregate_type, version, event_type, event_data)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: AppendToOutbox :exec
INSERT INTO outbox_events (id, aggregate_id, aggregate_type, version, event_type, event_data)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetOutboxEvents :many
SELECT * FROM outbox_events
ORDER BY created_at ASC
LIMIT $1;

-- name: DeleteOutboxEvent :exec
DELETE FROM outbox_events
WHERE id = $1;

-- name: GetEvents :many
SELECT version, event_type, event_data
FROM events
WHERE aggregate_id = $1
ORDER BY version ASC;
