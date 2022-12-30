package sliding_window_counter

import (
	"context"
	"github.com/ghnexpress/traefik-ratelimit/log"
)

func (r *repository) GetAllRequestCountCurrentWindow(ctx context.Context, ip string) (int, error) {
	log.Log("Start get request count of ip ", ip)
	allReqCount, err := r.GetRequestCountByIP(ctx, ip)
	if err != nil {
		return 0, err
	}
	log.Log("all req count", allReqCount)
	sumAllRequest := 0
	for _, requestCount := range allReqCount {
		sumAllRequest += requestCount
	}
	log.Log("sum of all req ", sumAllRequest)
	return sumAllRequest, nil
}
