package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic("AWS konfigürasyonu yüklenemedi")
	}

	S3Client = s3.NewFromConfig(cfg)
}

func GetS3BucketName() string {
	return os.Getenv("S3_BUCKET_NAME")
}
