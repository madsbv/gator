-- name: CreateFeed :one
INSERT INTO feeds ( id, created_at, updated_at, name, url, user_id )
VALUES (
gen_random_uuid(),
NOW(),
NOW(),
$1,
$2,
$3
) RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetAllFeedsWithUsernames :many
SELECT feeds.name, feeds.url, users.name as user_name
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;
