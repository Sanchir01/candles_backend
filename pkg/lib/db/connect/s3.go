package connect

import (
	"context"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"sync"

	"github.com/Sanchir01/candles_backend/internal/config"

	"golang.org/x/exp/slog"
)

type S3Storage struct {
	mu   sync.Mutex
	file map[string][]byte
}

func NewStorage() *S3Storage {
	return &S3Storage{
		file: make(map[string][]byte),
	}
}

func NewS3(ctx context.Context, lg *slog.Logger, cfg *config.Config) *s3.Client {
	creds := credentials.NewStaticCredentialsProvider(cfg.S3Store.Key, os.Getenv("S3_SECRET"), "")

	newawsconfig, err := awscfg.LoadDefaultConfig(
		ctx,
		awscfg.WithRegion(cfg.S3Store.Region),
		awscfg.WithCredentialsProvider(creds),
		awscfg.WithRegion(cfg.S3Store.Region),
	)

	if err != nil {
		lg.Error("ошибка при инициализации s3 хранилища", err.Error())
		return nil
	}
	awsS3Client := s3.NewFromConfig(newawsconfig)

	return awsS3Client
}
