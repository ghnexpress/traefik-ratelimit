package local_cache

import (
	"context"
)

const (
	MaxRetries = 10
)

func (r *localCacheRepository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	r.LocalCache.ReadModifyStore(ip, func(userRequestCount map[int]int, param map[string]interface{}) map[int]int {
		currPart := param["current_part"].(int)
		userRequestCount[currPart] += 1
		return userRequestCount
	}, map[string]interface{}{"current_part": part})

	return nil
}
