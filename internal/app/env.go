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
	Config       *config.Config
	Repositories *Repositories
	Services     *Services
}

func NewEnv() (*Env, error) {
	cfg := config.InitConfig()
	lg := setupLogger(cfg.Env)
	ctx := context.Background()

	pgxdb, err := NewDataBases(cfg)
	if err != nil {
		lg.Error("pgx error connect", err.Error())
		return nil, err
	}
	s3client, err := NewStorages(ctx, lg, cfg)
	if err != nil {
		lg.Error("s3 error connect", err.Error())
	}
	repos := NewRepositories(pgxdb)
	servises := NewServices(repos, s3client)
	env := Env{
		Logger:       lg,
		DataBase:     pgxdb,
		Storages:     s3client,
		Config:       cfg,
		Services:     servises,
		Repositories: repos,
	}

	return &env, nil
}
