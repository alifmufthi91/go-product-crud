package cache

import (
	"context"
	"encoding/json"
	"product-crud/util/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

var (
	redisClient *redis.Client
)

func InitCache(redis *redis.Client) {
	redisClient = redis
}

func Set(key string, value interface{}) {
	logger.Info("Set cache for key: %s", key)
	ctx := context.Background()
	err := redisClient.Set(ctx, key, value, time.Minute*5).Err()
	if err != nil {
		panic(err)
	}
	logger.Info("Cache is set for key: %s", key)
}

func Get(key string, result interface{}) *error {
	logger.Info("Get cache for key: %s", key)
	ctx := context.Background()
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return &err
	}
	logger.Info("Cache is get for key: %s", key)
	return nil
}
