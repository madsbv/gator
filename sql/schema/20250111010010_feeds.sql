-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
id UUID PRIMARY KEY,
created_at TIMESTAMP WITH TIME ZONE,
updated_at TIMESTAMP WITH TIME ZONE,
name TEXT,
url TEXT UNIQUE NOT NULL,
user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
