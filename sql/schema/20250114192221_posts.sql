-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
id UUID PRIMARY KEY,
created_at TIMESTAMP WITH TIME ZONE NOT NULL,
updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
title TEXT NOT NULL,
url TEXT UNIQUE NOT NULL,
description TEXT,
published_at TIMESTAMP WITH TIME ZONE,
feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
