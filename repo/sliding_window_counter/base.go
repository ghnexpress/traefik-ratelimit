package sliding_window_counter

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/sliding_window_counter"
)

type repository struct {
	memory *memcache.Client
}

func NewSlidingWindowCounterRepository(memory *memcache.Client) sliding_window_counter.Repository {
	return &repository{memory: memory}
}
