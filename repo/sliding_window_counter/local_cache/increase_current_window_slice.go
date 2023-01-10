package local_cache

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

const (
	MaxRetries = 10
)

func (r *localCacheRepository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()
	value, ok := r.LocalCache.Load(ip)
	if !ok {
		return fmt.Errorf("can't load data of ip %s %v", ip, value)
	}
	userRequestCount, ok := value.(map[int]int)
	if !ok {
		return fmt.Errorf("can't cast value %v type %s to map[int]int to load data of ip %s", userRequestCount, reflect.TypeOf(value), ip)
	}
	userRequestCount[part] += 1
	r.LocalCache.Store(ip, userRequestCount)
	return nil
}
