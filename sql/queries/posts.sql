-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT * 
FROM posts
JOIN feed_follows ff ON posts.feed_id = ff.feed_id
JOIN users u ON ff.user_id = u.id
WHERE u.name = $1
ORDER BY published_at DESC
LIMIT $2;

-- name: GetPostByUrl :one
SELECT * FROM posts WHERE url = $1;
