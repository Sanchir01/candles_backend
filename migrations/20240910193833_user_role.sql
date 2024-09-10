-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN admin TO role;


ALTER TABLE users
    ALTER COLUMN role TYPE VARCHAR(50);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
