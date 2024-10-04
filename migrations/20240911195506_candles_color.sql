-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='phone') THEN
ALTER TABLE users ADD COLUMN phone VARCHAR(11) NOT NULL UNIQUE;
END IF;


    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='candles' AND column_name='color_id') THEN
ALTER TABLE candles ADD COLUMN color_id UUID NOT NULL REFERENCES color(id);
END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


