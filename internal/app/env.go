package app

import (
	"context"
	"fmt"
	grpclientauth "github.com/Sanchir01/candles_backend/internal/server/grpc/auth"
	grpcclientorder "github.com/Sanchir01/candles_backend/internal/server/grpc/order"
	"log"
	"log/slog"

	"github.com/Sanchir01/candles_backend/internal/config"
)

type Env struct {
	Logger        *slog.Logger
	DataBase      *Database
	Storages      *Storages
	Config        *config.Config
	GRPCAuth      *grpclientauth.ClientAuthGRPC
	GRPCOrder     *grpcclientorder.ClientOrderGRPC
	Repositories  *Repositories
	Services      *Services
	KafkaProducer map[string]*Producer
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
	producers := make(map[string]*Producer)

	for _, topic := range cfg.Kafka.Outbox.Topic {
		kaf, err := NewProducer(cfg.Kafka.Outbox.Broke, topic, cfg.Kafka.Outbox.Retries, ctx)
		if err != nil {
			log.Fatalf("failed to create producer for topic %s: %v", topic, err)
		}
		producers[topic] = kaf
	}

	authgrpcurl := fmt.Sprintf("%s:%s", cfg.GrpcClients.GRPCAuth.Host, cfg.GrpcClients.GRPCAuth.Port)
	ordergrpcurl := fmt.Sprintf("%s:%s", cfg.GrpcClients.GRPCOrder.Host, cfg.GrpcClients.GRPCOrder.Port)
	authgrpc, err := grpclientauth.NewClientAuthGRPC(lg, authgrpcurl, cfg.GrpcClients.GRPCAuth.Timeout, cfg.GrpcClients.GRPCAuth.Retries)
	if err != nil {
		return nil, err
	}
	ordergrpc, err := grpcclientorder.NewClientOrderGRPC(lg, ordergrpcurl, cfg.GrpcClients.GRPCOrder.Timeout, cfg.GrpcClients.GRPCOrder.Retries)
	if err != nil {
		return nil, err
	}

	repos := NewRepositories(pgxdb)
	servises := NewServices(repos, s3client, pgxdb, producers["metrics"], lg)

	env := Env{
		Logger:        lg,
		DataBase:      pgxdb,
		Storages:      s3client,
		Config:        cfg,
		GRPCAuth:      authgrpc,
		GRPCOrder:     ordergrpc,
		Services:      servises,
		Repositories:  repos,
		KafkaProducer: producers,
	}

	return &env, nil
}
