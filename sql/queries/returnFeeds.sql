-- name: GetUserFeeds :many
select feeds.name, feeds.url, users.name AS username
FROM feeds
JOIN users ON feeds.user_id = users.id;