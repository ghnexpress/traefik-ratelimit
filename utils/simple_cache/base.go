package simple_local_cache

import "sync"

type Cache struct {
	sync.Map
}

func NewSimpleLocalCache() *Cache {
	return &Cache{}
}
