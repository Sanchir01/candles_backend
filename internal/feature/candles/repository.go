package candles

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}
type CandlesSortEnum int

const (
	CREATED_AT_ASC CandlesSortEnum = iota
	CREATED_AT_DESC
	SORT_RANK_ASC
	SORT_RANK_DESC
	PRICE_ASC
	PRICE_DESC
)

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{
		primaryDB,
	}
}
func (r *Repository) AllCandles(ctx context.Context, sort *model.CandlesSortEnum, categoryId, colorId *uuid.UUID) ([]model.Candles, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var orders string
	if sort != nil {
		orders = BuildSortQuery(*sort) // Разыменовываем указатель
	}
	slog.Error("filter", categoryId, colorId)
	queryBuilder := sq.
		Select("id, title, slug, price, images, version, category_id, created_at, updated_at").
		From("public.candles").
		Where(sq.Eq{"category_id": categoryId, "color_id": colorId}).
		PlaceholderFormat(sq.Dollar)

	if orders != "" {
		queryBuilder = queryBuilder.OrderBy(orders)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	slog.Error("query", query)
	// Выполняем запрос
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candles := make([]model.Candles, 0)

	for rows.Next() {
		var candle model.Candles
		if err := rows.Scan(&candle.ID, &candle.Title, &candle.Slug, &candle.Price, &candle.Images, &candle.Version, &candle.CategoryID, &candle.CreatedAt, &candle.UpdatedAt); err != nil {
			return nil, err
		}
		candles = append(candles, candle)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return candles, nil
}

func (r *Repository) CreateCandles(
	ctx context.Context, categoryID, colorID uuid.UUID, title string, slug string, images []string, price int, tr pgx.Tx,
) (uuid.UUID, error) {
	var id uuid.UUID
	query, arg, err := sq.
		Insert("public.candles").
		Columns("color_id", "title", "slug", "price", "images", "category_id").
		Values(colorID, title, slug, price, images, categoryID).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	if err := tr.QueryRow(ctx, query, arg...).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil

}

func (r *Repository) CandlesBySlug(ctx context.Context, slug string) (*model.Candles, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, args, err := sq.Select("id, title, slug, price, images, version, category_id, created_at, updated_at, color_id").
		From("public.candles").
		Where(sq.Eq{"slug": slug}).PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	//query := "SELECT id ,title,slug, price, images, version, category_id, created_at, updated_at, color_id FROM public.candles WHERE slug = $1"

	var candle DBCandles
	if err := conn.QueryRow(ctx, query, args...).Scan(
		&candle.ID,
		&candle.Title,
		&candle.Slug,
		&candle.Price,
		&candle.Images,
		&candle.Version,
		&candle.CategoryID,
		&candle.CreatedAt,
		&candle.UpdatedAt,
		&candle.ColorID,
	); err != nil {
		return nil, err
	}
	return (*model.Candles)(&candle), nil
}

func (r *Repository) CandlesById(ctx context.Context, id uuid.UUID) (*model.Candles, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.Select("id, title, slug, price, images, version, category_id, created_at, updated_at, color_id").
		From("public.candles").
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	fmt.Printf("Generated SQL: %s\n", query)
	fmt.Printf("Arguments: %v\n", args)
	var candle DBCandles
	if err := conn.QueryRow(ctx, query, args...).Scan(
		&candle.ID,
		&candle.Title,
		&candle.Slug,
		&candle.Price,
		&candle.Images,
		&candle.Version,
		&candle.CategoryID,
		&candle.CreatedAt,
		&candle.UpdatedAt,
		&candle.ColorID,
	); err != nil {
		slog.Error("error db candle", err.Error())
		return nil, err
	}

	return (*model.Candles)(&candle), nil
}

func (r *Repository) UpdateCandles(ctx context.Context, updates map[string]interface{}, id uuid.UUID) (string, error) {
	if len(updates) == 0 {
		return "", fmt.Errorf("no fields to update")
	}

	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return "", err
	}

	defer conn.Release()

	query := sq.
		Update("candles").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).Suffix("RETURNING id")
	for field, value := range updates {
		query = query.Set(field, value)
	}
	querybuilder, arg, err := query.ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build SQL query: %w", err)
	}

	var updatedID string
	err = conn.QueryRow(ctx, querybuilder, arg...).Scan(&updatedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no candle found with id: %s", id)
		}
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	// Возвращаем ID обновленной записи
	return updatedID, nil
}

func (r *Repository) DeleteCandlesById(ctx context.Context, id uuid.UUID) (string, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return "", err
	}

	defer conn.Release()

	query, args, err := sq.
		Delete("candles").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}
	var deletedID string
	err = conn.QueryRow(ctx, query, args...).Scan(&deletedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no candle found with id: %s", id)
		}
		return "", fmt.Errorf("failed to execute query: %w", err)
	}
	return deletedID, nil
}

func BuildSortQuery(sort model.CandlesSortEnum) string {
	switch sort {
	case model.CandlesSortEnumCreatedAtAsc:
		return "created_at ASC"
	case model.CandlesSortEnumCreatedAtDesc:
		return "created_at DESC"
	case model.CandlesSortEnumPriceAsc:
		return "price ASC"
	case model.CandlesSortEnumPriceDesc:
		return "price DESC"
	default:
		return ""
	}
}
