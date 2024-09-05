-- +goose Up
-- +goose StatementBegin
ALTER TABLE candles ADD COLUMN images TEXT[];
ALTER TABLE category DROP COLUMN images;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
