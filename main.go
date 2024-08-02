package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucketName = "aws-demo-test-bucket-912fi3"
const regionName = "eu-west-1"

func main() {
	var (
		s3Client *s3.Client
		err      error
		out      []byte
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s3Client, err = initS3Client(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "initS3Client error: %s\n", err)
		os.Exit(1)
	}
	if err = createS3Bucket(ctx, s3Client); err != nil {
		fmt.Fprintf(os.Stderr, "createS3Bucket error: %s\n", err)
		os.Exit(1)
	}
	if err = uploadToS3Bucket(ctx, manager.NewUploader(s3Client), "test.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "uploadToS3Bucket error: %s\n", err)
		os.Exit(1)
	}
	if out, err = downloadFromS3(ctx, manager.NewDownloader(s3Client)); err != nil {
		fmt.Fprintf(os.Stderr, "downloadFromS3 error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Download Complete: %s\n", out)
}
