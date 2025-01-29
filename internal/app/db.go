package app

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/Sanchir01/candles_backend/pkg/lib/db/connect"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	PrimaryDB *pgxpool.Pool
}

func NewDataBases(cfg *config.Config) (*Database, error) {
	pgxdb, err := connect.PGXNew(cfg, context.Background())
	if err != nil {
		return nil, err
	}

	return &Database{PrimaryDB: pgxdb}, nil
}

func (databases *Database) Close() error {
	databases.PrimaryDB.Close()
	return nil
}
