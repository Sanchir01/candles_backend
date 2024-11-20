package main

import (
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/app"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/go-faker/faker/v4"
	"log/slog"
)

func main() {
	env, err := app.NewEnv()

	candle, err := generateFakeCategory()
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
func generateFakeCategory() (*model.Category, error) {
	var candle *model.Category
	if err := faker.FakeData(&candle); err != nil {
		return nil, fmt.Errorf("failed to generate fake title: %w", err)
	}

	return candle, nil
}

func generateFakeColor() (*model.Color, error) {
	var candle *model.Color
	if err := faker.FakeData(&candle); err != nil {
		return nil, fmt.Errorf("failed to generate fake title: %w", err)
	}

	return candle, nil
}

func generateFakeUser() (*model.User, error) {
	var candle *model.User
	if err := faker.FakeData(&candle); err != nil {
		return nil, fmt.Errorf("failed to generate fake title: %w", err)
	}

	return candle, nil
}
