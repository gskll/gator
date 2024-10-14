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
SELECT posts.*
FROM feed_follows
INNER JOIN users on feed_follows.user_id = users.id
INNER JOIN posts on feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC NULLS LAST
LIMIT $2;
