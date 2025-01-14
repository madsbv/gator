-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds
DROP COLUMN last_fetched_at;
-- +goose StatementEnd
