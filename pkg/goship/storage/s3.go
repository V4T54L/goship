package storage

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client creates a new AWS S3 client using the provided credentials and region.
//
// Parameters:
//   - accessKey: AWS access key ID (e.g., "AKIA...").
//   - secretKey: AWS secret access key (e.g., "abcd1234...").
//   - region: AWS region (e.g., "us-east-1").
//
// Example:
//
//	client, err := NewS3Client("AKIA...", "abcd...", "us-east-1")
func NewS3Client(accessKey, secretKey, region string) (*s3.Client, error) {
	cfg := aws.Config{
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		Region:      region,
	}
	return s3.NewFromConfig(cfg), nil
}

// UploadFile uploads data to the specified S3 bucket and key.
func UploadFile(ctx context.Context, client *s3.Client, bucket, key string, file multipart.File) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	return err
}
