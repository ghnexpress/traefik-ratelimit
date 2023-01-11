package simple_local_cache

import "sync"

type Cache struct {
	mu   sync.RWMutex
	Data map[string]map[int]int
}

func NewSimpleLocalCache() *Cache {
	return &Cache{Data: map[string]map[int]int{}}
}

func (c *Cache) Store(key string, data map[int]int) {
	c.mu.Lock()
	c.Data[key] = data
	c.mu.Unlock()
}

func (c *Cache) Load(key string) (data map[int]int, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, ok = c.Data[key]
	return
}

func (c *Cache) ReadModifyStore(key string, modifyFn func(map[int]int, int) map[int]int, param int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data := c.Data[key]
	data = modifyFn(data, param)
	c.Data[key] = data
}
