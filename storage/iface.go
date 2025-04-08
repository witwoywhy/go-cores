package storage

import (
	"context"
	"time"
)

type Storager interface {
	Upload(ctx context.Context, request *UploadRequest) error
	GetUrl(filePath string) string
	GetSignedUrl(ctx context.Context, filePath string, ttl time.Duration) (string, error)
	Delete(ctx context.Context, filePath string) error
}
