-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
id UUID PRIMARY KEY,
created_at TIMESTAMP WITH TIME ZONE,
updated_at TIMESTAMP WITH TIME ZONE,
user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL,
UNIQUE (user_id, feed_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
