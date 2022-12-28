package sliding_window_counter

import "context"

func (r *repository) GetAllRequestCountCurrentWindow(ctx context.Context, ip string) (int, error) {
	allReqCount, err := r.GetRequestCountByIP(ctx, ip)
	if err != nil {
		return 0, err
	}
	sumAllRequest := 0
	for _, requestCount := range allReqCount {
		sumAllRequest += requestCount
	}
	return sumAllRequest, nil
}
