-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD  COLUMN phone VARCHAR(11) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
