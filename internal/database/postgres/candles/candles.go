package pgstorecandles

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type CandlesPostgresStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CandlesPostgresStore {
	return &CandlesPostgresStore{
		db: db,
	}
}

func (s *CandlesPostgresStore) AllCandles(ctx context.Context) ([]model.Candles, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var candles []dbCandles

	if err := conn.SelectContext(ctx, &candles, "SELECT * FROM candles"); err != nil {
		return nil, err
	}

	return lo.Map(candles, func(candles dbCandles, _ int) model.Candles { return model.Candles(candles) }), nil
}
func (s *CandlesPostgresStore) CreateCandles(ctx context.Context, categoryID uuid.UUID, title string, slug string) (uuid.UUID, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Close()
	var id uuid.UUID

	if err := conn.GetContext(ctx, &id, "INSERT INTO candles (category_id, title, slug) VALUES ($1, $2, $3) RETURNING id", categoryID, title, slug); err != nil {
		return uuid.Nil, err
	}
	return id, nil

}

func (s *CandlesPostgresStore) CandlesBySlug(ctx context.Context, slug string) (*model.Candles, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var candle dbCandles

	if err := conn.GetContext(ctx, &candle, "SELECT * FROM candles WHERE slug = $1", slug); err != nil {
		return nil, err
	}
	return (*model.Candles)(&candle), nil
}

type dbCandles struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"name"`
	Slug       string    `db:"slug"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Version    uint      `db:"version"`
	CategoryID uuid.UUID `db:"category_id"`
}
