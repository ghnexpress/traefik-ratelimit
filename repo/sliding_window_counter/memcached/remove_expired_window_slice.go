package memcached

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

func (r *memcachedRepository) RemoveExpiredWindowSlice(ctx context.Context, ip string, currSlice, windowTime int) (err error) {
	var data *memcache.Item
	allReqCount := make(map[int]int, 0)

	for i := 0; i < MaxRetries; i++ {
		if data, err = r.Memcached.Get(ip); err != nil {
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
		//r.ErrPub.SendError(fmt.Errorf("evicted part %v", partToBeEvicted))

		for _, part := range partToBeEvicted {
			delete(allReqCount, part)
		}

		if data.Value, err = json.Marshal(allReqCount); err != nil {
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
