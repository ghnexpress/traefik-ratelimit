package sliding_window_counter

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

func (r *repository) RemoveExpiredWindowSlice(ctx context.Context, ip string, currSlice, windowTime int) (err error) {
	var data *memcache.Item
	allReqCount := make(map[int]int, 0)

	for i := 0; i < MaxRetries; i++ {
		if data, err = r.Memory.Get(ip); err != nil {
			return err
		}

		err = json.Unmarshal(data.Value, &allReqCount)

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

		if data.Value, err = json.Marshal(allReqCount); err != nil {
			return err
		}

		if err = r.Memory.CompareAndSwap(data); err != nil {
			if err == memcache.ErrCASConflict {
				continue
			}
			return err
		}

	}

	return err
}
