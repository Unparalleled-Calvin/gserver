package server

import (
	"context"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisAddr := os.Getenv("GSERVER_REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "", // 无密码，留空字符串
			DB:       0,  // 使用默认 DB
		})
	})
	return rdb
}

func SetKey(key string, val interface{}) error {
	ctx := context.Background()
	rdb := GetRedisClient()
	return rdb.Set(ctx, key, val, 0).Err()
}

func GetKey(key string) (interface{}, error) {
	ctx := context.Background()
	rdb := GetRedisClient()
	return rdb.Get(ctx, key).Result()
}
