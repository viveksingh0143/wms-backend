package cache

import "sync"

type InMemoryCache struct {
	cache map[string]interface{}
	mutex sync.RWMutex
}

func (im *InMemoryCache) Init() error {
	im.cache = make(map[string]interface{})
	return nil
}

func (im *InMemoryCache) Get(key string) (interface{}, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	val, ok := im.cache[key]
	return val, ok
}

func (im *InMemoryCache) Set(key string, value interface{}) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	im.cache[key] = value
}

func (im *InMemoryCache) Delete(key string) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	delete(im.cache, key)
}
