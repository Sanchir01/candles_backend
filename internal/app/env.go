package app

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/config"
	"log/slog"
)

type Env struct {
	Logger       *slog.Logger
	DataBase     *Database
	Storages     *Storages
	Repositories *Repositories
	Services     *Services
}

func NewEnv(cfg *config.Config, ctx context.Context) (*Env, error) {
	lg := setupLogger(cfg.Env)

	pgxdb, err := NewDataBases(cfg)
	if err != nil {
		lg.Error("pgx error connect", err.Error())
		return nil, err
	}

	s3client, err := NewS3(ctx, lg, cfg)
	if err != nil {
		lg.Error("s3 error connect", err.Error())
	}
	env := Env{
		Logger:   lg,
		DataBase: pgxdb,
		Storages: s3client,
	}

	return &env, nil
}
