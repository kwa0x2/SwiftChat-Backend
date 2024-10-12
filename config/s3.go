package config

import (
	"context"
	"github.com/getsentry/sentry-go"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client is the global Amazon S3 client instance.
var S3Client *s3.Client

// region "InitS3" initializes the Amazon S3 client with the specified region from environment variables.
func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		sentry.CaptureException(err)
		panic("Failed to load AWS configuration") // Panic if there is an error while loading AWS configuration.
	}

	S3Client = s3.NewFromConfig(cfg) // Create a new S3 client from the loaded configuration.
}

// endregion

// region "GetS3BucketName" returns the name of the S3 bucket from environment variables.
func GetS3BucketName() string {
	return os.Getenv("S3_BUCKET_NAME")
}

// endregion
