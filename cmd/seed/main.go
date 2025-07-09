package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/candles_backend/internal/app"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func main() {
	env, err := app.NewEnv()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := env.DataBase.PrimaryDB.Acquire(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()
	colorId, err := env.Services.ColorService.CreateColor(ctx, gofakeit.Color())
	categoryId, err := env.Services.CategoryService.CreateCategory(ctx, gofakeit.Word())
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for i := 0; i < 100; i++ {

		if err != nil {
			slog.Error(err.Error())
			return
		}
		candlesID, err := env.Services.CandlesService.CreateCandles(
			ctx,
			categoryId,
			colorId,
			faker.Name(),
			faker.Word(),
			nil,
			100,
			500,
		)
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		if err != nil {
			slog.Error(err.Error())
			tx.Rollback(ctx)
			return
		}
		slog.Warn("candle id", candlesID)
	}
	if err := tx.Commit(ctx); err != nil {
		slog.Error(err.Error())
	}
}

func generateFakeCandle(colorId, categoryId uuid.UUID, tr pgx.Tx) (uuid.UUID, error) {
	var candle *model.Candles
	var candleId uuid.UUID
	if err := faker.FakeData(&candle); err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate fake title: %w", err)
	}
	name := gofakeit.Name()
	slug, err := utils.Slugify(name)
	if err != nil {
		return uuid.Nil, err
	}
	imagesArray := generateImageArrayUrl()

	query, arg, err := sq.
		Insert("public.candles").
		Columns("color_id", "title", "slug", "price", "images", "category_id", "description", "weight").
		Values(colorId, gofakeit.Word(), slug, gofakeit.IntRange(500, 5000), imagesArray, categoryId, gofakeit.LoremIpsumSentence(10), gofakeit.IntRange(100, 800)).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build SQL query: %w", err)
	}
	if err := tr.QueryRow(context.Background(), query, arg...).Scan(&candleId); err != nil {
		return uuid.Nil, err
	}
	return candleId, nil
}

func generateImageArrayUrl() []string {
	var urls []string
	for i := 0; i < 2; i++ {
		urls = append(urls, "https://random.imagecdn.app/800/600")
	}

	return urls
}
