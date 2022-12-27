package sliding_window_counter

import "context"

func (r *repository) GetAllRequestCountCurrentWindow(ctx context.Context, ip string) (int, error) {
	r.memory.Get(ip)
	return 0, nil
}
