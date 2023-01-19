package database

import (
	"fmt"
	"product-crud/config"
	"product-crud/util/logger"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisConnOnce sync.Once
	redisClient   *redis.Client
)

func RedisConnection() *redis.Client {
	redisConnOnce.Do(func() {
		var env = config.Env
		redisHost := env.RedisHost
		redisPort := env.RedisPort
		redisPassword := env.RedisPassword
		address := fmt.Sprintf("%s:%s", redisHost, redisPort)
		logger.Info("trying to connect Redis")
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: redisPassword,
			DB:       0,
		})

		redisClient = rdb
	})
	return redisClient
}
