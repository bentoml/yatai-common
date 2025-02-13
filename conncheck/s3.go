package conncheck

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type S3Probe struct {
	client *minio.Client
}

func NewS3Probe(c *minio.Client) *S3Probe {
	return &S3Probe{client: c}
}

func (p *S3Probe) Test(ctx context.Context, prefix string, bucketName string) error {
	objectName := fmt.Sprintf("%s-%d.txt", prefix, time.Now().UnixNano())
	content := []byte("test content")

	logrus.Infof("Testing S3 connection to %s", p.client.EndpointURL())
	exists, err := p.client.BucketExists(ctx, bucketName)
	if err != nil {
		logrus.Errorf("Failed to check bucket existence: %v", err)
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		logrus.Infof("Bucket %s does not exist, creating it", bucketName)
		err = p.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Errorf("Failed to create bucket: %v", err)
			return fmt.Errorf("failed to create test bucket: %w", err)
		}
		logrus.Infof("Bucket %s created successfully", bucketName)
	}

	reader := bytes.NewReader(content)
	_, err = p.client.PutObject(ctx, bucketName, objectName, reader, int64(len(content)), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		logrus.Errorf("Failed to write object: %v", err)
		return fmt.Errorf("failed to write object: %w", err)
	}
	logrus.Infof("Object %s written successfully", objectName)

	obj, err := p.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logrus.Errorf("Failed to get object: %v", err)
		return fmt.Errorf("failed to get object: %w", err)
	}
	logrus.Infof("Object %s retrieved successfully", objectName)
	defer obj.Close()

	readContent, err := io.ReadAll(obj)
	if err != nil {
		logrus.Errorf("Failed to read object content: %v", err)
		return fmt.Errorf("failed to read object content: %w", err)
	}
	logrus.Infof("Object content read successfully")

	if !bytes.Equal(readContent, content) {
		logrus.Errorf("Content mismatch: got %s, want %s", readContent, content)
		return fmt.Errorf("content mismatch: got %s, want %s", readContent, content)
	}

	err = p.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		logrus.Errorf("Failed to cleanup test object: %v", err)
		return fmt.Errorf("failed to cleanup test object: %w", err)
	}
	logrus.Infof("Test object %s cleaned up successfully", objectName)

	return nil
}
