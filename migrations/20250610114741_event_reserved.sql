-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN iF NOT EXISTS reserved_to TIMESTAMP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
