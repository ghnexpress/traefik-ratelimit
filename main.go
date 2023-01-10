package traefik_ratelimit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ghnexpress/traefik-ratelimit/utils"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/config"
	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter"
	slidingWindowCounter "github.com/ghnexpress/traefik-ratelimit/rate_limiter/sliding_window_counter"
	slidingWindowCounterLocalCache "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter/local_cache"
	slidingWindowCounterMemcached "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter/memcached"
	"github.com/ghnexpress/traefik-ratelimit/telegram"
	simple_local_cache "github.com/ghnexpress/traefik-ratelimit/utils/simple_cache"
)

const xRequestIDHeader = "X-Request-Id"

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *config.Config {
	return &config.Config{Memcached: config.MemcachedConfig{}}
}

type RateLimit struct {
	name                  string
	next                  http.Handler
	rate                  int
	memcachedRateLimiter  rate_limiter.RateLimiter
	localCacheRateLimiter rate_limiter.RateLimiter
	config                *config.Config
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *config.Config, name string) (http.Handler, error) {
	log.Log(fmt.Sprintf("config %v", config))
	memcachedInstance := memcache.New(config.Memcached.Address)
	memcachedInstance.Timeout = 500 * time.Millisecond
	localCache := simple_local_cache.NewSimpleLocalCache()

	telegramService := telegram.NewTelegramService(config.Telegram)
	slidingWindowCounterMemcachedRepository := slidingWindowCounterMemcached.NewSlidingWindowCounterMemcachedRepository(memcachedInstance, telegramService)
	memcachedRateLimiter := slidingWindowCounter.NewSlidingWindowCounter(
		slidingWindowCounterMemcachedRepository,
		telegramService,
		slidingWindowCounter.SlidingWindowCounterParam{
			MaxRequestInWindow: config.MaxRequestInWindow,
			WindowTime:         config.WindowTime,
		},
	)

	slidingWindowCounterLocalCachedRepository := slidingWindowCounterLocalCache.NewSlidingWindowCounterLocalCacheRepository(
		localCache,
	)
	localCacheRateLimiter := slidingWindowCounter.NewSlidingWindowCounter(
		slidingWindowCounterLocalCachedRepository,
		telegramService,
		slidingWindowCounter.SlidingWindowCounterParam{
			MaxRequestInWindow: config.MaxRequestInWindow,
			WindowTime:         config.WindowTime,
		})

	return &RateLimit{
		name:                  name,
		next:                  next,
		memcachedRateLimiter:  memcachedRateLimiter,
		localCacheRateLimiter: localCacheRateLimiter,
		config:                config,
	}, nil
}

func (r *RateLimit) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(rw)
	requestID := req.Header.Get(xRequestIDHeader)

	reqCtx := req.Context()
	reqCtx = context.WithValue(reqCtx, "requestID", requestID)
	reqCtx = context.WithValue(reqCtx, "env", r.config.Env)

	rw.Header().Add("Request-Ip", utils.GetIp(req))
	//if r.localCacheRateLimiter.IsAllowed(reqCtx, req, rw) {
	//	if r.memcachedRateLimiter.IsAllowed(reqCtx, req, rw) {
	//		r.next.ServeHTTP(rw, req)
	//		return
	//	}
	//}

	if r.memcachedRateLimiter.IsAllowed(reqCtx, req, rw) {
		r.next.ServeHTTP(rw, req)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusTooManyRequests)
	encoder.Encode(map[string]any{"status_code": http.StatusTooManyRequests, "message": "rate limit exceeded, try again later"})
	return
}
