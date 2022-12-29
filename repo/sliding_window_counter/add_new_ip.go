package sliding_window_counter

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/log"
)

func (r *repository) AddNewIP(ctx context.Context, ip string) error {
	body, err := json.Marshal(map[int]int{})
	if err != nil {
		return err
	}
	log.Log(body)
	return r.Memory.Set(&memcache.Item{
		Key:   ip,
		Value: body,
	})
}
