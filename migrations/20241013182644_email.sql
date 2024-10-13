-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN iF NOT EXISTS email TEXT NOT NULL;
ALTER TABLE users ADD COLUMN iF NOT EXISTS password TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
