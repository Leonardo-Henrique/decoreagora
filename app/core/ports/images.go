package ports

import (
	"context"
	"io"
	"sync"
	"time"
)

type ImagesHandler interface {
	SaveImage(ctx context.Context, file io.Reader, filename, bucketname string) error
	DeleteImage(ctx context.Context, wg *sync.WaitGroup, filename, bucketname string) error
	SignURL(ctx context.Context, bucket, key string, expireDuration time.Duration) (string, error)
}
