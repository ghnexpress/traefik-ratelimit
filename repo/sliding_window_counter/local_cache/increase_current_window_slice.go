package local_cache

import (
	"context"
)

const (
	MaxRetries = 10
)

func (r *localCacheRepository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	r.LocalCache.ReadModifyStore(ip, func(userRequestCount map[int]int, part int) map[int]int {
		userRequestCount[part] += 1
		return userRequestCount
	}, part)

	return nil
}
