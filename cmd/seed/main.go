package main

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/candles_backend/internal/app"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"math/rand"
	"time"
)

func randomColorName() string {
	rand.Seed(time.Now().UnixNano())
	colorNames := []string{
		"Red", "Green", "Blue", "Yellow", "Purple", "Cyan", "Magenta",
		"Orange", "Pink", "Brown", "Black", "White", "Gray", "Turquoise",
		"Gold", "Silver", "Ivory", "Coral", "Lavender",
	}
	return colorNames[rand.Intn(len(colorNames))]
}

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

	colorName := randomColorName()
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
	colorId, err := env.Services.ColorService.CreateColor(ctx, colorName)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	categoryId, err := env.Services.CategoryService.CreateCategory(ctx, faker.Word())
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
	slug, err := utils.Slugify(candle.Title)
	if err != nil {
		return uuid.Nil, err
	}
	imagesArray := generateImageArrayUrl()

	query, arg, err := sq.
		Insert("public.candles").
		Columns("color_id", "title", "slug", "price", "images", "category_id", "description", "weight").
		Values(colorId, candle.Title, slug, candle.Price, imagesArray, categoryId, candle.Description, candle.Weight).
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

func generateImageURL() string {
	width, err := faker.RandomInt(300, 1920)
	if err != nil {
		return ""
	}
	height, err := faker.RandomInt(300, 1080)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("https://picsum.photos/%d/%d", width, height)
}

func generateImageArrayUrl() []string {
	var urls []string
	for i := 0; i < 2; i++ {
		urls = append(urls, generateImageURL())
	}

	return urls
}
