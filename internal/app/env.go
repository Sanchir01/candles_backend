package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcserver "github.com/Sanchir01/candles_backend/internal/app/grpc"

	"github.com/Sanchir01/candles_backend/internal/config"
)

var address = []string{"localhost:9092"}

type Env struct {
	Logger        *slog.Logger
	DataBase      *Database
	Storages      *Storages
	Config        *config.Config
	GRPCSrv       *grpcserver.App
	Repositories  *Repositories
	Services      *Services
	KafkaProducer *Producer
}

func NewEnv() (*Env, error) {
	cfg := config.InitConfig()
	fmt.Println(cfg)
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
		return nil, err
	}
	kaf, err := NewProducer(cfg.Kafka.Producer.Broke, "order")
	if err != nil {
		lg.Error("kafka error connect", err.Error())
		return nil, err
	}
	repos := NewRepositories(pgxdb)
	servises := NewServices(repos, s3client)
	grpcApp := grpcserver.NewGrpc(lg, cfg.Grpc.Port)
	env := Env{
		Logger:        lg,
		DataBase:      pgxdb,
		Storages:      s3client,
		Config:        cfg,
		GRPCSrv:       grpcApp,
		Services:      servises,
		Repositories:  repos,
		KafkaProducer: kaf,
	}

	return &env, nil
}
