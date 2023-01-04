package local_cache

import (
	"github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
	simple_local_cache "github.com/ghnexpress/traefik-ratelimit/utils/simple_cache"
)

type localCacheRepository struct {
	LocalCache *simple_local_cache.Cache
}

func NewSlidingWindowCounterLocalCacheRepository(memory *simple_local_cache.Cache) sliding_window_counter.Repository {
	return &localCacheRepository{LocalCache: memory}
}
