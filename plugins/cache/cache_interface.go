package cache

type CacheInterface interface {
	Init() error
	Get(string) (interface{}, bool)
	Set(string, interface{})
	Delete(string)
}
