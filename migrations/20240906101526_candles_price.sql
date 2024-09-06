-- +goose Up
-- +goose StatementBegin
ALTER TABLE candles ADD COLUMN price NUMERIC DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
