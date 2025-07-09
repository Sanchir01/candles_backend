package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/Sanchir01/candles_backend/pkg/lib/db/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	PrimaryDB *pgxpool.Pool
	RedisDB   *redis.Client
}

func NewDataBases(cfg *config.Config, log *slog.Logger) (*Database, error) {
	pgxdb, err := connect.PGXNew(cfg, context.Background())
	if err != nil {
		log.Error("pgx connect error", err.Error())
		return nil, err
	}
	redisdb, err := connect.RedisConnect(context.Background(), cfg.RedisDB.Host, cfg.RedisDB.Port,
		os.Getenv("REDIS_PASSWORD"), cfg.Env,
		cfg.RedisDB.DBNumber, cfg.RedisDB.Retries)
	if err != nil {
		log.Error("redis connect error", err.Error())
		return nil, err
	}
	return &Database{PrimaryDB: pgxdb, RedisDB: redisdb}, nil
}

func (databases *Database) Close() error {
	databases.PrimaryDB.Close()

	return nil
}
