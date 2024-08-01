package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const bucketName = "aws-demo-test-bucket-912fi3"
const regionName = "eu-west-1"

func main() {
	var (
		s3Client *s3.Client
		err      error
		out      []byte
	)
	ctx := context.Background()

	if s3Client, err = initS3Client(ctx); err != nil {
		fmt.Printf("initS3Client error: %s", err)
		os.Exit(1)
	}
	if err = createS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("initS3Client error: %s", err)
		os.Exit(1)
	}
	if err = uploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("uploadToS3Bucket error: %s", err)
		os.Exit(1)
	}
	if out, err = downloadFromS3(ctx, s3Client); err != nil {
		fmt.Printf("DownloadFromS3 error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Download Complete: %s\n", out)
}

func initS3Client(ctx context.Context) (*s3.Client, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(regionName))
	if err != nil {
		return nil, fmt.Errorf("Unable to laod SDK config, %s", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("ListBucket error:, %s", err)
	}

	found := false
	for _, bucket := range allBuckets.Buckets {
		if *bucket.Name == bucketName {
			found = true
		}
	}
	if !found {
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: regionName,
			},
		})
		if err != nil {
			return fmt.Errorf("CreateBucket error:, %s", err)
		}
	}

	return nil
}

func uploadToS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	testFile, err := os.ReadFile("test.txt")
	if err != nil {
		return fmt.Errorf("Uplaod error:, %s", err)
	}
	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test.txt"),
		Body:   bytes.NewReader(testFile),
	})
	if err != nil {
		return fmt.Errorf("Uplaod error:, %s", err)
	}

	return nil
}

func downloadFromS3(ctx context.Context, s3Client *s3.Client) ([]byte, error) {
	downloader := manager.NewDownloader(s3Client)
	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test.txt"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file, %s", err)
	}

	if numBytesReceived := len(buffer.Bytes()); numBytes != int64(numBytesReceived) {
		return nil, fmt.Errorf("numBytes and numBytesReceived doesnt match %d %d", numBytes, numBytesReceived)
	}
	return buffer.Bytes(), nil
}
