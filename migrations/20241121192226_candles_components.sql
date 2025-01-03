-- +goose Up
-- +goose StatementBegin
CREATE TABLE candles_components (
                                    candle_id UUID NOT NULL,
                                    component_id UUID NOT NULL,
                                    PRIMARY KEY (candle_id, component_id),
                                    CONSTRAINT fk_candle FOREIGN KEY (candle_id) REFERENCES candles (id) ON DELETE CASCADE,
                                    CONSTRAINT fk_component FOREIGN KEY (component_id) REFERENCES components (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
