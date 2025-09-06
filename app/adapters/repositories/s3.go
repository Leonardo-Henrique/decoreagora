package repositories

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	awsConfig aws.Config
}

func NewS3Repository(awsConfig aws.Config) *S3Repository {
	return &S3Repository{
		awsConfig: awsConfig,
	}
}

func (s *S3Repository) newS3Client() *s3.Client {
	return s3.NewFromConfig(s.awsConfig)
}

func (s *S3Repository) SaveImage(ctx context.Context, file io.Reader, filename, bucketname string) error {
	s3client := s.newS3Client()
	uploader := manager.NewUploader(s3client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Repository) SignURL(ctx context.Context, bucket, key string, expireDuration time.Duration) (string, error) {
	s3client := s.newS3Client()
	presignClient := s3.NewPresignClient(s3client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expireDuration
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *S3Repository) DeleteImage(ctx context.Context, wg *sync.WaitGroup, filename, bucketname string) error {
	defer wg.Done()

	s3client := s.newS3Client()

	_, err := s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}

	// Wait until object is deleted
	waiter := s3.NewObjectNotExistsWaiter(s3client)
	err = waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(filename),
	}, 5*time.Minute) // 5 minute timeout
	if err != nil {
		return err
	}

	return nil
}
