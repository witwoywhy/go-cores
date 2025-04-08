package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient

	bucketName string
	publicUrl  string
	endPoint   string
	acl        string
}

func NewS3Storage(cfg *Config) Storager {
	s3Config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
	)

	if cfg.Region != "" {
		s3Config.Region = cfg.Region
	}

	if err != nil {
		panic(fmt.Errorf("failed to loaded default config on s3 storage: %v", err))
	}

	client := s3.NewFromConfig(s3Config, func(o *s3.Options) {
		if cfg.Target != "s3" {
			o.BaseEndpoint = &cfg.Endpoint
		}
	})

	return &s3Storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    cfg.BucketName,
		publicUrl:     cfg.PublicUrl,
		endPoint:      cfg.Endpoint,
		acl:           cfg.Acl,
	}
}

func (s *s3Storage) Upload(ctx context.Context, request *UploadRequest) error {
	i := &s3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &request.FilePath,
		Body:        request.File,
		ContentType: &request.ContentType,
	}
	if s.acl != "" {
		i.ACL = types.ObjectCannedACL(s.acl)
	}

	_, err := s.client.PutObject(ctx, i)

	return err
}

func (s *s3Storage) GetUrl(filePath string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicUrl, s.bucketName, filePath)
}

func (s *s3Storage) GetSignedUrl(ctx context.Context, filePath string, ttl time.Duration) (string, error) {
	signed, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &filePath,
	}, func(po *s3.PresignOptions) {
		po.Expires = ttl
	})
	if err != nil {
		return "", err
	}

	return signed.URL, nil
}

func (s *s3Storage) Delete(ctx context.Context, filePath string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    &filePath,
	})

	return err
}
