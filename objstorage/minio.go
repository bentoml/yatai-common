package objstorage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minioClient implements Client interface using MinIO
type minioClient struct {
	client *minio.Client
	cfg    Config
}

func newMinioClient(cfg Config) (Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	return &minioClient{
		client: client,
		cfg:    cfg,
	}, nil
}

// Upload implements Client.Upload for minioClient
func (c *minioClient) Upload(ctx context.Context, bucket, objectKey string, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := c.client.PutObject(ctx, bucket, objectKey, reader, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}
	return nil
}

// Download implements Client.Download for minioClient
func (c *minioClient) Download(ctx context.Context, bucket, objectKey string) ([]byte, error) {
	object, err := c.client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}
	return data, nil
}

// Delete implements Client.Delete for minioClient
func (c *minioClient) Delete(ctx context.Context, bucket, objectKey string) error {
	err := c.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// PresignedURL implements Client.PresignedURL for minioClient
func (c *minioClient) PresignedURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (*url.URL, error) {
	if c.cfg.Provider == AWS {
		url, err := c.client.PresignedPutObject(ctx, bucket, objectKey, expires)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
		}
		return url, nil
	}

	return c.gcsPresignURL(ctx, "PUT", bucket, objectKey, expires, nil, nil), nil
}

// GCSPresignURL returns a presigned URL for the given HTTP method, bucket, and object name.
// Implemented here because the GCS library doesn't support HMAC signing, and as of 2023-05-03,
// the minio library GCS presigned URL support is broken.
func (c *minioClient) gcsPresignURL(ctx context.Context, httpMethod string, bucketName string, objectName string, expires time.Duration, queryParameters url.Values, headers http.Header) *url.URL {
	escapedObjectName := strings.ReplaceAll(url.PathEscape(objectName), "%2F", "/")
	canonicalURI := "/" + escapedObjectName

	now := time.Now().UTC()
	requestTimestamp := now.Format("20060102T150405Z")

	dateStamp := now.Format("20060102")

	credentialScope := fmt.Sprintf("%s/%s/storage/goog4_request", dateStamp, c.cfg.Region)
	credential := fmt.Sprintf("%s/%s", c, credentialScope)

	host := fmt.Sprintf("%s.storage.googleapis.com", bucketName)

	if headers == nil {
		headers = http.Header{}
	}
	headers.Set("host", host)

	canonicalHeaders := ""
	sortedHeaderKeys := make([]string, 0, len(headers))
	for k := range headers {
		sortedHeaderKeys = append(sortedHeaderKeys, k)
	}
	sort.Strings(sortedHeaderKeys)
	for _, k := range sortedHeaderKeys {
		canonicalHeaders += fmt.Sprintf("%s:%s\n", strings.ToLower(k), strings.Join(headers[k], ","))
	}

	signedHeaders := ""
	for _, k := range sortedHeaderKeys {
		signedHeaders += strings.ToLower(k) + ";"
	}
	signedHeaders = signedHeaders[:len(signedHeaders)-1]

	if queryParameters == nil {
		queryParameters = url.Values{}
	}
	queryParameters.Add("x-goog-algorithm", "GOOG4-HMAC-SHA256")
	queryParameters.Add("x-goog-credential", credential)
	queryParameters.Add("x-goog-date", requestTimestamp)
	queryParameters.Add("x-goog-expires", strconv.Itoa(int(expires.Seconds())))
	queryParameters.Add("x-goog-signedheaders", signedHeaders)

	canonicalQueryString := queryParameters.Encode()

	canonicalRequest := strings.Join([]string{
		httpMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		"UNSIGNED-PAYLOAD",
	}, "\n")

	canonicalRequestSum := sha256.Sum256([]byte(canonicalRequest))

	strToSign := strings.Join([]string{
		"GOOG4-HMAC-SHA256",
		requestTimestamp,
		credentialScope,
		hex.EncodeToString(canonicalRequestSum[:]),
	}, "\n")

	keyDate := hmac.New(sha256.New, []byte("GOOG4"+c.cfg.SecretKey))
	keyDate.Write([]byte(dateStamp))
	keyDateSum := keyDate.Sum(nil)
	keyRegion := hmac.New(sha256.New, keyDateSum)
	keyRegion.Write([]byte(c.cfg.Region))
	keyRegionSum := keyRegion.Sum(nil)
	keyService := hmac.New(sha256.New, keyRegionSum)
	keyService.Write([]byte("storage"))
	keyServiceSum := keyService.Sum(nil)
	signingKey := hmac.New(sha256.New, keyServiceSum)
	signingKey.Write([]byte("goog4_request"))
	signingKeySum := signingKey.Sum(nil)
	messageDigest := hmac.New(sha256.New, signingKeySum)
	messageDigest.Write([]byte(strToSign))

	signature := hex.EncodeToString(messageDigest.Sum(nil))

	return &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     canonicalURI,
		RawQuery: canonicalQueryString + "&X-Goog-Signature=" + signature,
	}
}

// MakeBucket implements Client.MakeBucket for minioClient
func (c *minioClient) MakeBucket(ctx context.Context, bucket string) error {
	exists, err := c.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
			Region: c.cfg.Region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}
