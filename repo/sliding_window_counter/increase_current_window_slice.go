package sliding_window_counter

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MaxRetries = 10
)

func (r *repository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	for i := 0; i < MaxRetries; i++ {
		data, err := r.memory.Get(ip)
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
		if err = r.memory.CompareAndSwap(data); err != nil {
			if err == memcache.ErrCASConflict {
				continue
			}
			return err
		}
	}
	return err
}
