package candles

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{
		primaryDB,
	}
}
func (r *Repository) AllCandles(ctx context.Context) ([]model.Candles, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := "SELECT id ,title,slug, price, images, version, category_id, created_at, updated_at FROM public.candles"

	rows, err := conn.Query(ctx, query)
	if rows.Err(); err != nil {
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
	query := "INSERT INTO candles (category_id, title, slug, images,price,color_id) VALUES ($1, $2, $3, $4, $5,$6) RETURNING id"

	if err := tr.QueryRow(ctx, query, categoryID, title, slug, images, price, colorID).Scan(&id); err != nil {
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

	query := "SELECT id ,title,slug, price, images, version, category_id, created_at, updated_at, color_id FROM public.candles WHERE slug = $1"

	var candle DBCandles
	if err := conn.QueryRow(ctx, query, slug).Scan(
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

	query := "SELECT id ,title,slug, price, images, version, category_id, created_at, updated_at ,color_id FROM public.candles WHERE id = $1"

	var candle DBCandles
	if err := conn.QueryRow(ctx, query, id).Scan(
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