package Redis

import (
	"github.com/go-redis/redis/v8"
)

func GetRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	return rdb
}
