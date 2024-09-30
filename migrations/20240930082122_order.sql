-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id UUID NOT NULL REFERENCES users(id),
  status TEXT NOT NULL,
  total_amount NUMERIC NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  price NUMERIC NOT NULL,
  order_id UUID NOT NULL REFERENCES orders(id),
  quantity INT NOT NULL,
  product_id UUID NOT NULL REFERENCES candles(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
