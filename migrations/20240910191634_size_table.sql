-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS color(
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            title VARCHAR(100) NOT NULL ,
            slug VARCHAR(100) NOT NULL UNIQUE,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            version INT NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON "color"
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
