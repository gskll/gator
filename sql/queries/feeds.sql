-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, created_at, updated_at, name, url)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.*, users.name AS user_name
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.created_at DESC;
