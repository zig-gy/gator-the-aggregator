-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follow (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *)
SELECT
    iff.id,
    iff.created_at,
    iff.updated_at,
    iff.user_id,
    iff.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow iff
INNER JOIN feeds f
    ON iff.feed_id = f.id
INNER JOIN users u
    ON iff.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follow ff
INNER JOIN  feeds f
    ON ff.feed_id = f.id
INNER JOIN users u
    ON ff.user_id = u.id
WHERE u.id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follow ff
WHERE ff.user_id = $1
    AND ff.feed_id = $2;