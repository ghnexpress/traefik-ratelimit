package memcached

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"math/rand"
	"time"
)

func (r *memcachedRepository) RemoveExpiredWindowSlice(ctx context.Context, ip string, currSlice, windowTime int) (err error) {
	var data *memcache.Item
	allReqCount := make(map[int]int, 0)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < MaxRetries; i++ {
		duration := 100 * time.Millisecond
		if data, err = r.Memcached.Get(ip); err != nil {
			continue
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
