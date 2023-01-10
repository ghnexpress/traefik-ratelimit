package local_cache

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

func (r *localCacheRepository) RemoveExpiredWindowSlice(ctx context.Context, ip string, currSlice, windowTime int) (err error) {
	var mu sync.RWMutex
	mu.Lock()
	defer mu.Unlock()
	value, ok := r.LocalCache.Load(ip)
	if !ok {
		return fmt.Errorf("can't load data of ip %s %v", ip, value)
	}
	allReqCount, ok := value.(map[int]int)
	if !ok {
		return fmt.Errorf("can't cast value %v type %s to map[int]int to load data of ip %s", allReqCount, reflect.TypeOf(value), ip)
	}

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

	r.LocalCache.Store(ip, allReqCount)

	return nil
}
