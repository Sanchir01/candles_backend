package pgstorecandles

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type CandlesPostgresStore struct {
	pgxdb *pgxpool.Pool
}

func New(pgxdb *pgxpool.Pool) *CandlesPostgresStore {
	return &CandlesPostgresStore{
		pgxdb: pgxdb,
	}
}
