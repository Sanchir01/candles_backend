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

func (s3str *S3Store) PutObjects(ctx context.Context, folderpath string, images []*graphql.Upload) ([]string, error) {
	imageURLs := make([]string, 0)
	for _, image := range images {
		if image.File == nil {
			return nil, nil
		}
		filekey := fmt.Sprintf("%s/%s", folderpath, image.Filename)
		_, err := s3str.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(s3str.cfg.S3Store.BucketName),
			Key:           aws.String(filekey),
			ContentLength: &image.Size,
			ContentType:   aws.String(image.ContentType),
			ACL:           types.ObjectCannedACLPublicRead,
			Body:          image.File,
		})
		if err != nil {
			return nil, err
		}
		imageURL := fmt.Sprintf("%s%s%s", s3str.cfg.S3Store.URL, s3str.cfg.S3Store.BucketName, "/"+filekey)
		imageURLs = append(imageURLs, imageURL)
	}
	return imageURLs, nil
}

func (s3str *S3Store) DeleteObjects(ctx context.Context, folderpath string, images []*graphql.Upload) error {
	for _, image := range images {
		if folderpath == "" {
			return fmt.Errorf("folderpath is empty")
		}
		filekey := fmt.Sprintf("%s/%s", folderpath, image.Filename)
		_, err := s3str.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s3str.cfg.S3Store.BucketName),
			Key:    aws.String(filekey),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
