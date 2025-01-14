-- name: CreateFeedFollow :one
WITH inserted_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
gen_random_uuid(),
NOW(),
NOW(),
$1,
$2
) RETURNING * )
SELECT inserted_follow.*, feeds.name as feed_name, users.name as user_name
FROM inserted_follow
INNER JOIN users
ON users.id = inserted_follow.user_id
INNER JOIN feeds
ON feeds.id = inserted_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name as feed_name, users.name as user_name
FROM feed_follows
INNER JOIN users
ON users.id = feed_follows.user_id
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollowByUrl :exec
WITH feed AS (
SELECT * FROM feeds
WHERE feeds.url = $1
)
DELETE FROM feed_follows
USING feed
WHERE feed_follows.feed_id = feed.id AND feed_follows.user_id = $2;
