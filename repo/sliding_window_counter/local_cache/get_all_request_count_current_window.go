package local_cache

import (
	"context"
	"github.com/ghnexpress/traefik-ratelimit/log"
)

func (r *localCacheRepository) GetAllRequestCountCurrentWindow(ctx context.Context, ip string) (int, error) {
	log.Log("Start get request count from local cache of ip ", ip)
	allReqCount, err := r.GetRequestCountByIP(ctx, ip)
	if err != nil {
		return 0, err
	}
	log.Log("all req count map from local cache", allReqCount)
	sumAllRequest := 0
	for _, requestCount := range allReqCount {
		sumAllRequest += requestCount
	}
	log.Log("sum of all req from local cache ", sumAllRequest)
	return sumAllRequest, nil
}
