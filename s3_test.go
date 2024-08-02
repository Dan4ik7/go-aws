package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MockS3Client struct {
	ListBucketsOutput  *s3.ListBucketsOutput
	CreateBucketOutput *s3.CreateBucketOutput
}

func (m MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsOutput, nil
}

func (m MockS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.CreateBucketOutput, nil
}

type MockS3Uploader struct {
	UploadOuput *manager.UploadOutput
}

func (m MockS3Uploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	return m.UploadOuput, nil
}

func TestCreateS3Bucket(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := createS3Bucket(ctx, MockS3Client{
		ListBucketsOutput: &s3.ListBucketsOutput{
			Buckets: []types.Bucket{
				{
					Name: aws.String("test-bucket"),
				},
				{
					Name: aws.String("test-bucket-2"),
				},
			},
		},
		CreateBucketOutput: &s3.CreateBucketOutput{},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "createS3Bucket error: %s\n", err)
		os.Exit(1)
	}
}

func TestUploadToS3Bucket(t *testing.T) {
	mockUploader := MockS3Uploader{
		UploadOuput: &manager.UploadOutput{},
	}
	err := uploadToS3Bucket(context.Background(), mockUploader, "testdata/test.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "UploadToS3Bucket error: %s\n", err)
		os.Exit(1)
	}
}
