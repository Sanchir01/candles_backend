-- +goose Up
-- +goose StatementBegin
ALTER TABLE size RENAME TO color;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
