package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	server string
}

func (rc *RedisCache) Init() error {
	rc.client = redis.NewClient(&redis.Options{
		Addr: rc.server,
	})
	return rc.client.Ping(context.Background()).Err()
}

func (rc *RedisCache) Get(key string) (interface{}, bool) {
	val, err := rc.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (rc *RedisCache) Set(key string, value interface{}) {
	rc.client.Set(context.Background(), key, value, 0)
}

func (rc *RedisCache) Delete(key string) {
	rc.client.Del(context.Background(), key)
}
