-- name: CreatePost :exec
INSERT INTO posts (id, feed_id, created_at, updated_at, title, url, description, published_at)
VALUES (
    $1,
    $2,
    NOW(),
    NOW(),
    $3,
    $4,
    $5,
    $6
);

-- name: GetPostsForUser :many
SELECT p.*, f.name AS feed_name
FROM posts p
INNER JOIN feeds f ON p.feed_id = f.id
INNER JOIN feed_follows ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC NULLS LAST
LIMIT $2;
