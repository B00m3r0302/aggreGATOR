-- name: GetUserId :one
SELECT id
FROM users
WHERE name = $1;