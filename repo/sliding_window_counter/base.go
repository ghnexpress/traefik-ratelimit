package sliding_window_counter

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
)

type Repository interface {
	GetRequestCountByIP(ctx context.Context, ip string) (map[int]int, error)
	IncreaseCurrentWindowSlice(ctx context.Context, key string, part int) error
	GetAllRequestCountCurrentWindow(ctx context.Context, key string) (int, error)
	AddNewIP(ctx context.Context, ip string) error
	RemoveExpiredWindowSlice(ctx context.Context, key string, currSlice, windowTime int) error
}

type repository struct {
	Memory *memcache.Client
}

func NewSlidingWindowCounterRepository(memory *memcache.Client) Repository {
	return &repository{Memory: memory}
}
