package local_cache

import (
	"context"
)

func (r *localCacheRepository) RemoveExpiredWindowSlice(ctx context.Context, ip string, currSlice, windowTime int) (err error) {
	r.LocalCache.ReadModifyStore(ip, func(allReqCount map[int]int, param map[string]interface{}) map[int]int {
		partToBeEvicted := make([]int, 0)
		for part, _ := range allReqCount {
			if part < currSlice-windowTime {
				partToBeEvicted = append(partToBeEvicted, part)
			}
		}

		if len(partToBeEvicted) < 0 {
			return nil
		}

		for _, part := range partToBeEvicted {
			delete(allReqCount, part)
		}
		return allReqCount
	}, nil)
	return nil
}
