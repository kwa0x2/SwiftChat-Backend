package config

import (
	"github.com/getsentry/sentry-go"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

// region "RedisSession" initializes a Redis session store with specific options.
func RedisSession() redis.Store {
	// Create a new Redis store with the given parameters.
	store, err := redis.NewStore(10, "tcp", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("SESSION_SECRET_KEY")))
	if err != nil {
		sentry.CaptureException(err)
		panic(err) // Panic if there is an error while creating the store.
	}

	// Set options for the session store.
	store.Options(sessions.Options{
		MaxAge:   86400, // Session duration in seconds (24 hours).
		Path:     "/",
		HttpOnly: true, // HTTP-only flag to prevent client-side scripts from accessing cookies.
		Secure:   true, // Set to true to ensure cookies are sent only over HTTPS.
	})

	return store // Return the configured Redis store.
}

// endregion
