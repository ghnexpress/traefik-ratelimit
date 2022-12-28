package sliding_window_counter

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

func (r *repository) GetRequestCountByIP(ctx context.Context, ip string) (requestCountPerIP map[int]int, err error) {
	var data *memcache.Item
	if data, err = r.Memory.Get(ip); err != nil {
		return nil, err
	}
	err = json.Unmarshal(data.Value, &requestCountPerIP)
	return requestCountPerIP, err
}
