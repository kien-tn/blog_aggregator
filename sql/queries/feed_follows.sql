-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, user_id, feed_id, created_at, updated_at)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) 
SELECT 
    ff.id,
    ff.user_id,
    ff.feed_id,
    ff.created_at,
    ff.updated_at,
    f.name AS feed_name,
    u.name AS user_name
FROM
    inserted_feed_follow ff
    JOIN feeds f ON ff.feed_id = f.id
    JOIN users u ON ff.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT 
    ff.id,
    f.name AS feed_name,
    u.name AS user_name
FROM
    feed_follows ff
    JOIN feeds f ON ff.feed_id = f.id
    JOIN users u ON ff.user_id = u.id
WHERE
    u.name = $1;

-- name: DropFeedFollowsForUrlCurrentUser :exec
DELETE FROM feed_follows
WHERE
    feed_id = (SELECT id FROM feeds WHERE url = $1)
    AND user_id = (SELECT id FROM users WHERE users.name = $2);