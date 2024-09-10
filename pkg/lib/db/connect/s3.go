package connect

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
)

func NewS3(ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						URL:           "http://localhost:4566",
						SigningRegion: "us-east-1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested for service %s in region %s", service, region)
			}),
		))
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}
	client := s3.NewFromConfig(cfg)
	bucketName := "test-bucket"
	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: "us-east-1",
		},
	})
	if err != nil {
		log.Fatalf("ошибка создания бакета: %v", err)
	}
	fmt.Printf("Бакет %s создан\n", bucketName)
}
