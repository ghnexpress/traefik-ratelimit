package memcached

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"math/rand"
	"time"
)

const (
	MaxRetries = 20
)

func (r *memcachedRepository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < MaxRetries; i++ {
		var data *memcache.Item
		duration := 100 * time.Millisecond
		data, err = r.Memcached.Get(ip)
		if err != nil {
			continue
		}
		userRequestCount := map[int]int{}
		if err = json.Unmarshal(data.Value, &userRequestCount); err != nil {
			return err
		}
		userRequestCount[part] += 1
		if data.Value, err = json.Marshal(userRequestCount); err != nil {
			return err
		}
		if err = r.Memcached.CompareAndSwap(data); err != nil {
			if err == memcache.ErrCASConflict {
				time.Sleep(duration)
				duration = time.Duration(float64(duration) * (rand.Float64() + 1))
				continue
			}
			return fmt.Errorf("compare and swap %v", err)
		}
		break
	}
	return err
}
