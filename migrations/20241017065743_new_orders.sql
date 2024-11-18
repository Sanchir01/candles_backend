-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
                                     id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                     status TEXT NOT NULL,
                                     total_amount NUMERIC NOT NULL,
                                     version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS order_items (
                                           id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           price NUMERIC NOT NULL,
                                           order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                                           quantity INT NOT NULL,
                                           product_id UUID NOT NULL REFERENCES candles(id) ON DELETE CASCADE,
                                           version INTEGER NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON "orders"
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
