package storage

import (
	"context"
	"sync"

	"github.com/Unparalleled-Calvin/gserver/internal/settings"
	"github.com/redis/go-redis/v9"
)

var (
	rdb       *redis.Client
	redisOnce sync.Once
)

func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     settings.RedisAddr,
			Password: settings.RedisPassword,
			DB:       settings.RedisDB,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			panic("Failed to connect to Redis: " + err.Error())
		}
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
