package main

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/candles_backend/internal/app"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

func main() {
	env, err := app.NewEnv()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ctx := context.Background()

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
	if err != nil {
		slog.Error(err.Error())
		return
	}

	categoryId, err := env.Services.CategoryService.CreateCategory(ctx, gofakeit.Word())
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Error("category id", categoryId, "color id", colorId)
	for i := 0; i < 10; i++ {
		candleId, err := generateFakeCandle(colorId, categoryId, tx)
		if err != nil {
			slog.Error(err.Error())
			tx.Rollback(ctx)
			return
		}
		slog.Warn("candle id", candleId)
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