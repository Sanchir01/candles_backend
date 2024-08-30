-- +goose Up
-- +goose StatementBegin
ALTER TABLE category ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();
ALTER TABLE category RENAME COLUMN new_id TO id;
UPDATE category SET new_id = uuid_generate_v4() WHERE new_id IS NULL;
ALTER TABLE category DROP COLUMN id;
ALTER TABLE category RENAME COLUMN new_id TO id;
ALTER TABLE category ADD PRIMARY KEY (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
