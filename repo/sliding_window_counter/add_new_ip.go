package sliding_window_counter

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

func (r *repository) AddNewIP(ctx context.Context, ip string) error {
	body, _ := json.Marshal(map[int]int{})
	return r.Memory.Set(&memcache.Item{
		Key:   ip,
		Value: body,
	})
}
