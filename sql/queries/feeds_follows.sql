-- name: CreateFeedFollow :one

WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (created_at,updated_at,user_id,feed_id) 
    VALUES
    (   $1,
        $2,
        $3,
        $4)
    RETURNING *
)
SELECT i.* ,f.name, u.name 
FROM inserted_feed_follow i JOIN feeds f ON i.feed_id=f.id JOIN users u ON i.user_id=u.id;

-- name: GetFeedFollowsForUser :many
SELECT f.name FROM feed_follows ff JOIN feeds f ON ff.feed_id=f.id JOIN users u ON ff.user_id=u.id where u.name=$1;

-- name: GetFeedURL :one
SELECT * FROM feeds WHERE url=$1;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows ff WHERE ff.feed_id=(SELECT id FROM feeds where url=$1) AND ff.user_id=$2
RETURNING *;
