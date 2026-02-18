-- name: QueryById :one
SELECT feeds.name AS feed_name,
       feeds.id,
       users.name AS user_name,
       feeds.url
FROM feeds
INNER JOIN users ON feeds.user_id = users.id
WHERE feeds.url = $1;