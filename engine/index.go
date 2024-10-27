package engine

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisHelper struct {
	*redis.Client
}

var redisHelper *RedisHelper

var redisOnce sync.Once

// GetRedisHelper 用于获取已经初始化的 RedisHelper 实例
func GetRedisHelper() *RedisHelper {
	redisOnce.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			PoolSize:     10,
			PoolTimeout:  30 * time.Second,
		})
		redisHelper = &RedisHelper{Client: rdb}
	})
	return redisHelper
}
