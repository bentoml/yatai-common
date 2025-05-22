package objstorage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// azureClient implements Client interface using Azure Blob Storage
type azureClient struct {
	client     *azblob.ServiceURL
	credential *azblob.SharedKeyCredential
}

func newAzureClient(cfg Config) (Client, error) {
	credential, err := azblob.NewSharedKeyCredential(cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create azure credentials: %w", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	url, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}
	serviceURL := azblob.NewServiceURL(*url, pipeline)

	return &azureClient{
		client:     &serviceURL,
		credential: credential,
	}, nil
}

// Upload implements Client.Upload for azureClient
func (c *azureClient) Upload(ctx context.Context, bucket, objectKey string, data []byte) error {
	containerURL := c.client.NewContainerURL(bucket)
	blobURL := containerURL.NewBlockBlobURL(objectKey)

	_, err := azblob.UploadBufferToBlockBlob(ctx, data, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload blob: %w", err)
	}
	return nil
}

// Download implements Client.Download for azureClient
func (c *azureClient) Download(ctx context.Context, bucket, objectKey string) ([]byte, error) {
	containerURL := c.client.NewContainerURL(bucket)
	blobURL := containerURL.NewBlockBlobURL(objectKey)

	response, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download blob: %w", err)
	}

	bodyStream := response.Body(azblob.RetryReaderOptions{})
	defer bodyStream.Close()

	data, err := io.ReadAll(bodyStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read blob data: %w", err)
	}
	return data, nil
}

// Delete implements Client.Delete for azureClient
func (c *azureClient) Delete(ctx context.Context, bucket, objectKey string) error {
	containerURL := c.client.NewContainerURL(bucket)
	blobURL := containerURL.NewBlockBlobURL(objectKey)

	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	return nil
}

// PresignedURL implements Client.PresignedURL for azureClient
func (c *azureClient) PresignedURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (*url.URL, error) {
	containerURL := c.client.NewContainerURL(bucket)
	blobURL := containerURL.NewBlockBlobURL(objectKey)

	permissions := azblob.BlobSASPermissions{Create: true, Write: true}
	start := time.Now().UTC()
	expiry := start.Add(expires)

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     start,
		ExpiryTime:    expiry,
		Permissions:   permissions.String(),
		ContainerName: bucket,
		BlobName:      objectKey,
	}.NewSASQueryParameters(c.credential)

	if err != nil {
		return nil, fmt.Errorf("failed to generate SAS query parameters: %w", err)
	}

	urlToSign := blobURL.URL()
	urlToSign.RawQuery = sasQueryParams.Encode()
	return &urlToSign, nil
}

// MakeBucket implements Client.MakeBucket for azureClient
func (c *azureClient) MakeBucket(ctx context.Context, bucket string) error {
	containerURL := c.client.NewContainerURL(bucket)

	_, err := containerURL.GetProperties(ctx, azblob.LeaseAccessConditions{})
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok && serr.ServiceCode() == azblob.ServiceCodeContainerNotFound {
			_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
			if err != nil {
				return fmt.Errorf("failed to create container: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check container existence: %w", err)
		}
	}
	return nil
}
