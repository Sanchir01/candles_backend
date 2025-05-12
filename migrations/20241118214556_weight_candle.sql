-- +goose Up
-- +goose StatementBegin
ALTER TABLE candles ADD COLUMN iF NOT EXISTS description text NOT NULL DEFAULT '';
ALTER TABLE candles ADD COLUMN iF NOT EXISTS weight numeric NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
