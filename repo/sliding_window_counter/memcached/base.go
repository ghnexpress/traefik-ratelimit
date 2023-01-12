package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
)

const (
	MaxRetries = 20
)

type memcachedRepository struct {
	Memcached *memcache.Client
	ErrPub    log.ErrorPublisher
}

func NewSlidingWindowCounterMemcachedRepository(memory *memcache.Client, errPublisher log.ErrorPublisher) sliding_window_counter.Repository {
	return &memcachedRepository{Memcached: memory, ErrPub: errPublisher}
}
