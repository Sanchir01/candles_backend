package connect

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"

	"log/slog"
)

func NewS3(ctx context.Context, lg *slog.Logger, cfg *config.Config) *s3.Client {
	creds := credentials.NewStaticCredentialsProvider(cfg.S3Store.Key, os.Getenv("S3_SECRET"), "")
	newawsconfig, err := awscfg.LoadDefaultConfig(
		ctx,
		awscfg.WithRegion(cfg.S3Store.Region),
		awscfg.WithCredentialsProvider(creds),
		awscfg.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: cfg.S3Store.URL, // URL для MinIO или LocalStack
				}, nil
			}),
		))

	if err != nil {
		lg.Error("ошибка при инициализации s3 хранилища", err.Error())
		return nil
	}
	awsS3Client := s3.NewFromConfig(newawsconfig)

	return awsS3Client
}
