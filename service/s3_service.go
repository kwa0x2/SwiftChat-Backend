package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kwa0x2/realtime-chat-backend/config"
)

type S3Service struct{}

func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)

	params := &s3.PutObjectInput{
		Bucket: aws.String(config.GetS3BucketName()),
		Key:    aws.String(fileName),
		Body:   file,
	}

	_, err := config.S3Client.PutObject(context.TODO(), params)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.GetS3BucketName(), fileName)
	return fileURL, nil
}
