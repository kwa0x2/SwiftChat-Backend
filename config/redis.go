package config

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func RedisConnection() sessions.Store {
    store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }

    fmt.Println("Redis connection successfuly")
    return store
}
