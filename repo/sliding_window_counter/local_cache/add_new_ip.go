package local_cache

import (
	"context"
)

func (r *localCacheRepository) AddNewIP(ctx context.Context, ip string) error {
	emptyMap := map[int]int{}
	r.LocalCache.Store(ip, emptyMap)
	return nil
}
