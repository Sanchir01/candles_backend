-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Добавление столбца uuid для категорий
ALTER TABLE category ADD COLUMN IF NOT EXISTS id UUID DEFAULT uuid_generate_v4();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
