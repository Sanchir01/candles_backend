package s3store

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Store struct {
	s3Client *s3.Client
	cfg      *config.Config
}

func New(s3Client *s3.Client, ctx context.Context, cfg *config.Config) *S3Store {
	return &S3Store{
		s3Client: s3Client,
		cfg:      cfg,
	}
}

func (s3str *S3Store) PutObjects(ctx context.Context, images []*graphql.Upload) ([]string, error) {
	imageURLs := make([]string, 0)
	for _, image := range images {
		if image.File == nil {
			return nil, nil
		}
		image, err := s3str.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(s3str.cfg.S3Store.BucketName),
			Key:           aws.String(image.Filename),
			ContentLength: &image.Size,
			ContentType:   aws.String(image.ContentType),
			ACL:           types.ObjectCannedACLPublicRead,
			Body:          image.File,
		})
		if err != nil {
			return nil, err
		}
		imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", "your-bucket-name", "your-region", fileKey)
		imageURLs = append(imageURLs, imageURL)
	}
	return imageURLs, nil
}
