-- GetUrlById :one
SELECT url
FROM feeds
WHERE id = $1;