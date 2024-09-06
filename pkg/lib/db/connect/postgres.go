package connect

import (
	"context"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"time"
)

func PostgresCon(cfg *config.Config, lg *slog.Logger) *sqlx.DB {
	postgresString := fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s password=%s port=%s host=%s",
		cfg.DB.User, cfg.DB.Database, cfg.DB.SSL, os.Getenv("PASSWORD_POSTGRES"),
		cfg.DB.Port, cfg.DB.Host,
	)

	db, err := sqlx.Open("postgres", postgresString)
	if err != nil {
		lg.Error("sqlx.Connect error", slog.String("error", err.Error()))
	}
	return db
}

func PGXNew(cfg *config.Config, lg *slog.Logger, ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DB.User, os.Getenv("PASSWORD_POSTGRES"),
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Database,
	)
	var pool *pgxpool.Pool
	var err error

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Создаем пул соединений
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
