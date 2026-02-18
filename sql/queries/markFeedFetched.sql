-- name: MarkLastFeedFetched :one
UPDATE feeds
SET last_fetched_at = sqlc.arg(last_fetched_at)::timestamptz,
    updated_at      = sqlc.arg(last_fetched_at)::timestamptz
WHERE id = sqlc.arg(id)
    RETURNING id, last_fetched_at, updated_at;
