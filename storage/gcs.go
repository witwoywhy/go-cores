package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type gcsStorage struct {
	client     *storage.BucketHandle
	bucketName string
	publicUrl  string
	prefix     string
}

func NewGcsStorage(cfg *Config) Storager {
	creds := option.WithCredentialsJSON([]byte(cfg.ServiceAccount))
	client, err := storage.NewClient(context.Background(), creds)
	if err != nil {
		panic(fmt.Errorf("failed to create gcs client: %v", err))
	}

	return &gcsStorage{
		client:     client.Bucket(cfg.BucketName),
		bucketName: cfg.BucketName,
		publicUrl:  cfg.PublicUrl,
		prefix:     cfg.Prefix,
	}
}

func (g *gcsStorage) Upload(ctx context.Context, request *UploadRequest) error {
	writer := g.client.Object(g.prefix + "/" + request.FilePath).NewWriter(ctx)
	writer.ContentType = request.ContentType
	defer writer.Close()

	if request.DownloadFileName != nil {
		writer.ObjectAttrs.ContentDisposition = fmt.Sprintf(`inline; filename="%v"`, *request.DownloadFileName)
	}

	_, err := io.Copy(writer, request.File)
	return err
}

func (g *gcsStorage) GetUrl(filePath string) string {
	return fmt.Sprintf("%s/%s/%s", g.publicUrl, g.prefix, filePath)
}

func (g *gcsStorage) GetSignedUrl(ctx context.Context, filePath string, ttl time.Duration) (string, error) {
	url, err := g.client.SignedURL(g.prefix+"/"+filePath, &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  http.MethodGet,
		Expires: time.Now().Add(ttl),
	})
	if err != nil {
		return "", err
	}

	return url, nil
}

func (g *gcsStorage) Delete(ctx context.Context, filePath string) error {
	return g.client.Object(g.prefix + "/" + filePath).Delete(ctx)
}
