package main

import (
	"encoding/gob"
	"github.com/getsentry/sentry-go"
	"github.com/kwa0x2/swiftchat-backend/internal/app"
	"log"
	"time"
)

func init() {
	gob.Register(time.Time{})
}

func main() {
	application := app.NewApp()
	application.SetupRoutes()

	if err := application.Run(); err != nil {
		sentry.CaptureException(err)
		log.Fatal("failed to run app: ", err)
	}
}
