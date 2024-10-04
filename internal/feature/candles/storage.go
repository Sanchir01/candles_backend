package candles

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	client          *s3.Client
	bucketName, url string
}

func NewStorage(client *s3.Client, bucketName, url string) *Storage {
	return &Storage{
		client,
		bucketName,
		url,
	}
}

func (s3str *Storage) PutObjects(ctx context.Context, folderpath string, images []*graphql.Upload) ([]string, error) {
	imageURLs := make([]string, 0)
	for _, image := range images {
		if image.File == nil {
			return nil, nil
		}
		filekey := fmt.Sprintf("/%s/%s", folderpath, image.Filename)
		_, err := s3str.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(s3str.bucketName),
			Key:           aws.String(filekey),
			ContentLength: &image.Size,
			ContentType:   aws.String(image.ContentType),
			ACL:           types.ObjectCannedACLPublicRead,
			Body:          image.File,
		})
		if err != nil {
			return nil, err
		}
		imageURL := fmt.Sprintf("%s%s%s", s3str.url, s3str.bucketName, "/"+filekey)
		imageURLs = append(imageURLs, imageURL)
	}
	return imageURLs, nil
}

func (s3str *Storage) DeleteObjects(ctx context.Context, folderpath string, images []*graphql.Upload) error {
	for _, image := range images {
		if folderpath == "" {
			return fmt.Errorf("folderpath is empty")
		}
		filekey := fmt.Sprintf("%s/%s", folderpath, image.Filename)
		_, err := s3str.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s3str.bucketName),
			Key:    aws.String(filekey),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
