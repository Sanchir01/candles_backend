package connect

import (
	"context"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

func PGXNew(cfg *config.Config, ctx context.Context, env string) (*pgxpool.Pool, error) {
	var dsn string
	switch env {
	case "development":
		dsn = fmt.Sprintf(
			"postgresql://postgres:postgres@localhost:5435/test",
		)
	case "production":
		dsn = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s",
			cfg.DB.User, os.Getenv("PASSWORD_POSTGRES"),
			cfg.DB.Host, cfg.DB.Port, cfg.DB.Database,
		)
	}

	var pool *pgxpool.Pool
	var err error

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, cfg.DB.MaxAttempts, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
