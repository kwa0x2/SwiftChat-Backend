package config

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func RedisSession() redis.Store {
	store, err := redis.NewStore(10, "tcp", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("SESSION_SECRET_KEY")))
	if err != nil {
		panic(err)
	}

	store.Options(sessions.Options{
		MaxAge: 86400,
		Path: "/",
		HttpOnly: true,
		Secure: true,
	})

    return store
}

