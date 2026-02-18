-- name: GetNextFeedToFetch :many
SELECT id,
       name,
       url,
       last_fetched_at,
       updated_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
