-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUEs (
gen_random_uuid(),
NOW(),
NOW(),
$1,
$2,
$3,
$4,
$5
) RETURNING *;

-- name: GetPostsForUser :many
SELECT * FROM posts
INNER JOIN
feed_follows
ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at
LIMIT $2;
