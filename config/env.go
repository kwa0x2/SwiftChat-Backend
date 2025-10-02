package config

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	PostgreUser         string `mapstructure:"POSTGRE_USER" validate:"required"`
	PostgrePassword     string `mapstructure:"POSTGRE_PASSWORD" validate:"required"`
	PostgreHost         string `mapstructure:"POSTGRE_HOST" validate:"required"`
	PostgreDB           string `mapstructure:"POSTGRE_DB" validate:"required"`
	RedirectURL         string `mapstructure:"REDIRECT_URL" validate:"required,url"`
	ClientID            string `mapstructure:"CLIENT_ID" validate:"required"`
	ClientSecret        string `mapstructure:"CLIENT_SECRET" validate:"required"`
	SentryDSN           string `mapstructure:"SENTRY_DSN" validate:"required"`
	AWSRegion           string `mapstructure:"AWS_REGION" validate:"required"`
	S3BucketName        string `mapstructure:"S3_BUCKET_NAME" validate:"required"`
	RedisHost           string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD" validate:"required"`
	SessionSecretKey    string `mapstructure:"SESSION_SECRET_KEY" validate:"required"`
	ResendAPIKey        string `mapstructure:"RESEND_API_KEY" validate:"required"`
	JWTSecretKey        string `mapstructure:"JWT_SECRET_KEY" validate:"required"`
}

var Env Config

func LoadEnv() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Error reading .env file: %v", err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&Env); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Failed to parse environment: %v", err)
		os.Exit(1)
	}

	validate := validator.New()
	if err := validate.Struct(Env); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Environment validation failed: %v", err)
		os.Exit(1)
	}
}
