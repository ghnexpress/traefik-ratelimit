package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
)

type memcachedRepository struct {
	Memcached *memcache.Client
}

func NewSlidingWindowCounterMemcachedRepository(memory *memcache.Client) sliding_window_counter.Repository {
	return &memcachedRepository{Memcached: memory}
}
