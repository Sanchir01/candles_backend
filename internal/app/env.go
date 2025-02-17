package app

import (
	"context"
	grpcserver "github.com/Sanchir01/candles_backend/internal/app/grpc"
	"log/slog"

	"github.com/Sanchir01/candles_backend/internal/config"
)

type Env struct {
	Logger       *slog.Logger
	DataBase     *Database
	Storages     *Storages
	Config       *config.Config
	GRPCSrv      *grpcserver.App
	Repositories *Repositories
	Services     *Services
	Bot          *Bot
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
	grpcApp := grpcserver.NewGrpc(lg, cfg.Grpc.Port)
	bot := NewBot(servises)
	env := Env{
		Logger:       lg,
		DataBase:     pgxdb,
		Storages:     s3client,
		Config:       cfg,
		GRPCSrv:      grpcApp,
		Services:     servises,
		Repositories: repos,
		Bot:          bot,
	}

	return &env, nil
}
