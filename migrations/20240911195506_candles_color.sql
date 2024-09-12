-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD  COLUMN phone VARCHAR(11) NOT NULL UNIQUE;

ALTER TABLE candles ADD COLUMN color_id UUID NOT NULL REFERENCES color(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
