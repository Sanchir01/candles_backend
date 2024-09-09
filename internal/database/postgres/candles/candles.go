package pgstorecandles

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"time"
)

type CandlesPostgresStore struct {
	db    *sqlx.DB
	pgxdb *pgxpool.Pool
}

func New(db *sqlx.DB, pgxdb *pgxpool.Pool) *CandlesPostgresStore {
	return &CandlesPostgresStore{
		db: db, pgxdb: pgxdb,
	}
}

func (s *CandlesPostgresStore) AllCandles(ctx context.Context) ([]model.Candles, error) {
	conn, err := s.pgxdb.Acquire(ctx)
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
	for {
		var candle model.Candles
		if err := rows.Scan(&candle.ID, &candle.Title, &candle.Slug, &candle.Price, &candle.Images, &candle.Version, &candle.CategoryID, &candle.CreatedAt, &candle.UpdatedAt); err != nil {
			return nil, err
		}
	}
	return candles, nil

}

func (s *CandlesPostgresStore) CreateCandles(
	ctx context.Context, categoryID uuid.UUID, title string, slug string, images []string,
) (uuid.UUID, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()
	var id uuid.UUID
	query := "INSERT INTO candles (category_id, title, slug, images) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := conn.QueryRow(ctx, query, categoryID, title, slug, images).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil

}

func (s *CandlesPostgresStore) CandlesBySlug(ctx context.Context, slug string) (*model.Candles, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id ,title,slug, price, images, version, category_id, created_at, updated_at FROM public.candles WHERE slug = $1"

	var candle dbCandles
	if err := conn.QueryRow(ctx, query, slug).Scan(&candle.ID, &candle.Title, candle.Slug, &candle.CreatedAt, &candle.UpdatedAt, &candle.Version); err != nil {
		return nil, err
	}
	return (*model.Candles)(&candle), nil
}

func (s *CandlesPostgresStore) CandlesById(ctx context.Context, id uuid.UUID) (*model.Candles, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id ,title,slug, price, images, version, category_id, created_at, updated_at FROM public.candles WHERE id = $1"

	var candle dbCandles
	if err := conn.QueryRow(ctx, query, id).Scan(&candle.ID, &candle.Title, candle.Slug, &candle.CreatedAt, &candle.UpdatedAt, &candle.Version); err != nil {
		return nil, err
	}
	return (*model.Candles)(&candle), nil
}

type dbCandles struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	Slug       string    `db:"slug"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Version    uint      `db:"version"`
	Price      int       `db:"price"`
	Images     []string  `db:"images"`
	CategoryID uuid.UUID `db:"category_id"`
}
