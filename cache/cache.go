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

func Set(ctx context.Context, key string, value interface{}) error {
	logger.Info("Set cache for key: %s", key)
	err := redisClient.Set(ctx, key, value, time.Minute*5).Err()
	if err != nil {
		return err
	}
	logger.Info("Cache is set for key: %s", key)
	return nil
}

func Get(ctx context.Context, key string, result interface{}) error {
	logger.Info("Get cache for key: %s", key)
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return err
	}
	logger.Info("Cache is get for key: %s", key)
	return nil
}
