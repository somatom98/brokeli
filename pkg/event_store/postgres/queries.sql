-- name: GetEventsByAggregateID :many
SELECT id, aggregate_type, aggregate_id, version, event_type, event_data, created_at
FROM events 
WHERE aggregate_type = $1 AND aggregate_id = $2
ORDER BY version ASC;

-- name: InsertEvent :execresult
INSERT INTO events (id, aggregate_type, aggregate_id, version, event_type, event_data)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetLatestVersionForAggregate :one
SELECT COALESCE(MAX(version), 0) as latest_version
FROM events 
WHERE aggregate_type = $1 AND aggregate_id = $2;

-- name: GetAllEvents :many
SELECT id, aggregate_type, aggregate_id, version, event_type, event_data, created_at
FROM events 
ORDER BY created_at ASC;