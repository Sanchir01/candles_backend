-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN role SET DEFAULT 'user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
