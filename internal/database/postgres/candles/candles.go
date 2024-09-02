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

type dbCandles struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"name"`
	Slug       string    `db:"slug"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Version    uint      `db:"version"`
	CategoryID uuid.UUID `db:"category_id"`
}
