package storage

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Target     string `mapstructure:"target"`
	BucketName string `mapstructure:"bucketName"`
	Prefix     string `mapstructure:"prefix"`
	PublicUrl  string `mapstructure:"publicUrl"`
	Region     string `mapstructure:"region"`
	S3Config   `mapstructure:",squash"`
	BlobConfig `mapstructure:",squash"`
	GcsConfig  `mapstructure:",squash"`
}

type S3Config struct {
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Endpoint  string `mapstructure:"endpoint"`
	Acl       string `mapstructure:"acl"`
}

type BlobConfig struct {
	AzureConnection string `mapstructure:"azureConnection"`
}

type GcsConfig struct {
	ServiceAccount string `mapstructure:"serviceAccount"`
}

func Init(key string) Storager {
	var config Config
	if err := viper.UnmarshalKey(key, &config); err != nil {
		panic(fmt.Errorf("failed to loaded storage config %s: %v", key, err))
	}

	switch config.Target {
	case "s3", "r2", "minio":
		return NewS3Storage(&config)
	case "blob":
		return NewAzureStorage(&config)
	case "gcs":
		return NewGcsStorage(&config)
	default:
		return nil
	}
}
