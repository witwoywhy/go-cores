package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

type azureStorage struct {
	client        *container.Client
	serviceClient *azblob.Client
	publicUrl     string
	prefix        string
}

func NewAzureStorage(cfg *Config) Storager {
	serviceClient, err := azblob.NewClientFromConnectionString(cfg.AzureConnection, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create service client azure blob: %v", err))
	}

	return &azureStorage{
		client:        serviceClient.ServiceClient().NewContainerClient(cfg.Prefix),
		serviceClient: serviceClient,
		publicUrl:     cfg.PublicUrl,
		prefix:        cfg.Prefix,
	}
}

func (a *azureStorage) Upload(ctx context.Context, request *UploadRequest) error {
	bc := a.client.NewBlockBlobClient(request.FilePath)
	_, err := bc.UploadStream(ctx, request.File, nil)
	return err
}

func (a *azureStorage) GetUrl(filePath string) string {
	return fmt.Sprintf("%s/%s/%s", a.publicUrl, a.prefix, filePath)
}

func (a *azureStorage) GetSignedUrl(ctx context.Context, filePath string, ttl time.Duration) (string, error) {
	bc := a.client.NewBlockBlobClient(filePath)
	now := time.Now()
	url, err := bc.GetSASURL(
		sas.BlobPermissions{Read: true},
		now.Add(ttl),
		&blob.GetSASURLOptions{StartTime: &now})
	if err != nil {
		return "", err
	}

	return url, nil
}

func (a *azureStorage) Delete(ctx context.Context, filePath string) error {
	_, err := a.serviceClient.DeleteBlob(ctx, a.prefix, filePath, nil)
	return err
}
