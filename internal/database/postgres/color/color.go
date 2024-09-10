package pgstorecolor

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ColorPostgresStore struct {
	pgxdb *pgxpool.Pool
}

func New(pgxdb *pgxpool.Pool) *ColorPostgresStore {
	return &ColorPostgresStore{
		pgxdb: pgxdb,
	}
}

func (s *ColorPostgresStore) SizeBySlug(ctx context.Context, slug string) (*model.Color, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.color WHERE slug = $1"

}
