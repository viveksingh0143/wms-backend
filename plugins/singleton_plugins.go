package plugins

import (
	cachePlugin "star-wms/plugins/cache"
	"sync"
)

// Singleton cache instance and sync.Once object
var (
	Cache cachePlugin.CacheInterface
	once  sync.Once
)

func GetCache() cachePlugin.CacheInterface {
	once.Do(func() {
		Cache = cachePlugin.CreateInMemoryCacheInstance()
	})
	return Cache
}
