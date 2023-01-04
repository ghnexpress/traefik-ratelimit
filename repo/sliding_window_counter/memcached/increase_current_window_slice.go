package memcached

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MaxRetries = 10
)

func (r *memcachedRepository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	for i := 0; i < MaxRetries; i++ {
		data, err := r.Memcached.Get(ip)
		if err != nil {
			return err
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
				continue
			}
			return err
		}
		break
	}
	return err
}
