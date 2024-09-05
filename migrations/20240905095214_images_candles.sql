-- +goose Up
-- +goose StatementBegin
ALTER TABLE category ADD COLUMN images TEXT[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd