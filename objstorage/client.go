package objstorage

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Provider represents supported cloud storage providers
type Provider string

const (
	AWS   Provider = "aws"
	GCP   Provider = "gcp"
	Azure Provider = "azure"
)

// Client interface defines common operations for object storage
type Client interface {
	Upload(ctx context.Context, bucket, objectKey string, data []byte) error
	Download(ctx context.Context, bucket, objectKey string) ([]byte, error)
	Delete(ctx context.Context, bucket, objectKey string) error
	PresignedURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (*url.URL, error)
	MakeBucket(ctx context.Context, bucket string) error
}

// Config holds the configuration for object storage client
type Config struct {
	Provider  Provider
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	UseSSL    bool
}

// NewClient creates a new object storage client based on the provider
func NewClient(cfg Config) (Client, error) {
	switch cfg.Provider {
	case AWS, GCP:
		return newMinioClient(cfg)
	case Azure:
		return newAzureClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
