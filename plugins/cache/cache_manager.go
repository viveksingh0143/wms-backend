package cache

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"sync"
)

type Type int

const (
	NoCache Type = iota
	InMemory
	Redis
)

type Manager struct {
	cacheType Type
}

var (
	once                 sync.Once
	cacheInstance        CacheInterface
	cacheManagerInstance *Manager
	errInit              error
)

func NewCacheManager() (*Manager, error) {
	once.Do(func() {
		var cacheType Type
		cacheTypeStr := viper.GetString("cache.type")

		switch cacheTypeStr {
		case "in_memory":
			cacheType = InMemory
		case "redis":
			cacheType = Redis
		case "no_cache":
			cacheType = NoCache
		default:
			errInit = fmt.Errorf("invalid cache type %s", cacheTypeStr)
			return
		}

		switch cacheType {
		case InMemory:
			cacheInstance = &InMemoryCache{}
		case Redis:
			cacheInstance = &RedisCache{server: viper.GetString("cache.redis.server")}
		case NoCache:
			cacheInstance = &NoOpCache{}
		}

		errInit = cacheInstance.Init()
		cacheManagerInstance = &Manager{cacheType: cacheType}
	})

	return cacheManagerInstance, errInit
}

func (cm *Manager) Get(key string) (interface{}, bool) {
	return cacheInstance.Get(key)
}

func (cm *Manager) Set(key string, value interface{}) {
	cacheInstance.Set(key, value)
}

func (cm *Manager) Delete(key string) {
	cacheInstance.Delete(key)
}

func (cm *Manager) GetOrCreate(key string, createFunc func() interface{}) (interface{}, bool) {
	if value, exists := cm.Get(key); exists {
		return value, true
	}

	newValue := createFunc()
	cm.Set(key, newValue)
	return newValue, false
}

func (cm *Manager) GetTarget(key string, targetType interface{}) (interface{}, bool) {
	rawValue, exists := cacheInstance.Get(key)
	if !exists {
		return targetType, false
	}

	targetTypeOf := reflect.TypeOf(targetType)
	if reflect.TypeOf(rawValue) == targetTypeOf {
		return rawValue, true
	}

	return nil, false
}

func (cm *Manager) GetOrCreateTarget(key string, createFunc func() interface{}, targetType interface{}) (interface{}, bool) {
	rawValue, exists := cm.GetTarget(key, targetType)
	if exists {
		return rawValue, true
	}

	newValue := createFunc()
	cm.Set(key, newValue)

	return cm.GetTarget(key, targetType)
}
