package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没密码留空
		DB:       0,  // 默认库
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
