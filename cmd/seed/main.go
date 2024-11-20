package main

import (
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/go-faker/faker/v4"
	"log/slog"
)

func main() {
	candle, err := generateFakeCandle()
	if err != nil {
		panic(err)
	}
	slog.Warn("candle", candle)
}
func generateFakeCandle() (*model.Candles, error) {
	var candle *model.Candles
	if err := faker.FakeData(&candle); err != nil {
		return nil, fmt.Errorf("failed to generate fake title: %w", err)
	}

	return candle, nil
}
