-- +goose Up
-- +goose StatementBegin
CREATE  TABLE IF NOT EXISTS components (
                                       id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                       title VARCHAR(100) NOT NULL ,
                                       slug VARCHAR(100) NOT NULL UNIQUE,
                                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       version INT NOT NULL DEFAULT 1
);
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_category_updated_at BEFORE UPDATE
    ON components FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();


ALTER TABLE candles ADD COLUMN iF NOT EXISTS description text NOT NULL DEFAULT '';
ALTER TABLE candles ADD COLUMN iF NOT EXISTS weight numeric NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
