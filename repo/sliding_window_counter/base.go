package sliding_window_counter

import (
	"context"
)

type Repository interface {
	GetRequestCountByIP(ctx context.Context, ip string) (map[int]int, error)
	IncreaseCurrentWindowSlice(ctx context.Context, key string, part int) error
	GetAllRequestCountCurrentWindow(ctx context.Context, key string) (int, error)
	AddNewIP(ctx context.Context, ip string) error
	RemoveExpiredWindowSlice(ctx context.Context, key string, currSlice, windowTime int) error
}
