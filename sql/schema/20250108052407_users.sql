-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
       id UUID PRIMARY KEY,
       created_at TIMESTAMP WITH TIME ZONE NOT NULL,
       updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
       name TEXT UNIQUE NOT NULL
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
