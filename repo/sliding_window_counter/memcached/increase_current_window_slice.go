package memcached

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MaxRetries = 20
)

func (r *memcachedRepository) IncreaseCurrentWindowSlice(ctx context.Context, ip string, part int) (err error) {
	for i := 0; i < MaxRetries; i++ {
		var data *memcache.Item
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
		//if err = r.Memcached.CompareAndSwap(data); err != nil {
		//	if err == memcache.ErrCASConflict {
		//		time.Sleep(100 * time.Millisecond)
		//		continue
		//	}
		//	return err
		//}
		if err = r.Memcached.Set(data); err != nil {
			//if err == memcache.ErrCASConflict {
			//	time.Sleep(100 * time.Millisecond)
			//	continue
			//}
			return err
		}
		break
	}
	return err
}
