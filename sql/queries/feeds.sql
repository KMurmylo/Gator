-- name: InsertFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT f.id, f.name AS feedName, f.url, u.name AS userName FROM feeds f JOIN users u ON u.id=f.user_id;
