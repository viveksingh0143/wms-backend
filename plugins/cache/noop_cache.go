package cache

type NoOpCache struct{}

func (nc *NoOpCache) Init() error {
	return nil
}

func (nc *NoOpCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func (nc *NoOpCache) Set(key string, value interface{}) {
	// Do nothing
}

func (nc *NoOpCache) Delete(key string) {
	// Do nothing
}
