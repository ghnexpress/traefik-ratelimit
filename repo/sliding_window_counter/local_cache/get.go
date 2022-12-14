package local_cache

import (
	"context"
	"fmt"
)

func (r *localCacheRepository) GetRequestCountByIP(ctx context.Context, ip string) (map[int]int, error) {
	requestCountPerIP, ok := r.LocalCache.Load(ip)
	if !ok {
		return nil, fmt.Errorf("can't load data for ip %s", ip)
	}
	if !ok {
		return nil, fmt.Errorf("can't load data for ip %s", ip)
	}
	return requestCountPerIP, nil
}
