package config

import (
	"github.com/getsentry/sentry-go"
	"log"
	"os"
)

func InitSentry() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %w", err)
	}
}