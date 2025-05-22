package app

import (
	"context"
	"fmt"
	grpclientauth "github.com/Sanchir01/candles_backend/internal/server/grpc"
	"log/slog"

	"github.com/Sanchir01/candles_backend/internal/config"
)

type Env struct {
	Logger        *slog.Logger
	DataBase      *Database
	Storages      *Storages
	Config        *config.Config
	GRPCAuth      *grpclientauth.ClientAuthGRPC
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
	topic := cfg.Kafka.Producer.Topic
	kaf, err := NewProducer(cfg.Kafka.Producer.Broke, topic)
	if err != nil {
		lg.Error("kafka error connect", err.Error())
		return nil, err
	}

	authgrpc, err := grpclientauth.NewClientAuthGRPC(lg, "0.0.0.0:44049", cfg.GrpcClients.GRPCAuth.Timeout, cfg.GrpcClients.GRPCAuth.Retries)
	if err != nil {
		return nil, err
	}
	repos := NewRepositories(pgxdb)
	servises := NewServices(repos, s3client)

	env := Env{
		Logger:        lg,
		DataBase:      pgxdb,
		Storages:      s3client,
		Config:        cfg,
		GRPCAuth:      authgrpc,
		Services:      servises,
		Repositories:  repos,
		KafkaProducer: kaf,
	}

	return &env, nil
}
