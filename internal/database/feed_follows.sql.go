// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
gen_random_uuid(),
NOW(),
NOW(),
$1,
$2
) RETURNING id, created_at, updated_at, user_id, feed_id )
SELECT inserted_follow.id, inserted_follow.created_at, inserted_follow.updated_at, inserted_follow.user_id, inserted_follow.feed_id, feeds.name as feed_name, users.name as user_name
FROM inserted_follow
INNER JOIN users
ON users.id = inserted_follow.user_id
INNER JOIN feeds
ON feeds.id = inserted_follow.feed_id
`

type CreateFeedFollowParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  sql.NullString
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow, arg.UserID, arg.FeedID)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const deleteFeedFollowByUrl = `-- name: DeleteFeedFollowByUrl :exec
WITH feed AS (
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
WHERE feeds.url = $1
)
DELETE FROM feed_follows
USING feed
WHERE feed_follows.feed_id = feed.id AND feed_follows.user_id = $2
`

type DeleteFeedFollowByUrlParams struct {
	Url    string
	UserID uuid.UUID
}

func (q *Queries) DeleteFeedFollowByUrl(ctx context.Context, arg DeleteFeedFollowByUrlParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowByUrl, arg.Url, arg.UserID)
	return err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, feed_follows.feed_id, feeds.name as feed_name, users.name as user_name
FROM feed_follows
INNER JOIN users
ON users.id = feed_follows.user_id
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
`

type GetFeedFollowsForUserRow struct {
	ID        uuid.UUID
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  sql.NullString
	UserName  string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.FeedName,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
