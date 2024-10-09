-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN iF NOT EXISTS version INTEGER NOT NULL  DEFAULT 1;

ALTER TABLE order_items ADD COLUMN iF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
