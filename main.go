package traefik_ratelimit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter/sliding_window_counter"
	slidingWindowCounterLocalCache "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter/local_cache"
	slidingWindowCounterMemcached "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter/memcached"
	simple_local_cache "github.com/ghnexpress/traefik-ratelimit/utils/simple_cache"
	"net/http"
)

type MemcachedConfig struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	MaxRequestInWindow int             `json:"maxRequestInWindow,omitempty"`
	WindowTime         int             `json:"windowTime,omitempty"`
	MemcachedConfig    MemcachedConfig `json:"memcachedConfig"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{MemcachedConfig: MemcachedConfig{}}
}

type RateLimit struct {
	name                  string
	next                  http.Handler
	rate                  int
	memcachedRateLimiter  rate_limiter.RateLimiter
	localCacheRateLimiter rate_limiter.RateLimiter
	config                *Config
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log.Log(fmt.Sprintf("config %v", config))
	memcachedInstance := memcache.New(config.MemcachedConfig.Address)
	localCache := simple_local_cache.NewSimpleLocalCache()
	log.Log(localCache)

	slidingWindowCounterMemcachedRepository := slidingWindowCounterMemcached.NewSlidingWindowCounterMemcachedRepository(memcachedInstance)
	memcachedRateLimiter := sliding_window_counter.NewSlidingWindowCounter(
		slidingWindowCounterMemcachedRepository,
		sliding_window_counter.SlidingWindowCounterParam{
			MaxRequestInWindow: config.MaxRequestInWindow,
			WindowTime:         config.WindowTime,
		},
	)

	slidingWindowCounterLocalCachedRepository := slidingWindowCounterLocalCache.NewSlidingWindowCounterLocalCacheRepository(
		localCache,
	)
	localCacheRateLimiter := sliding_window_counter.NewSlidingWindowCounter(
		slidingWindowCounterLocalCachedRepository,
		sliding_window_counter.SlidingWindowCounterParam{
			MaxRequestInWindow: config.MaxRequestInWindow,
			WindowTime:         config.WindowTime,
		})

	log.Log(fmt.Sprintf("%v", slidingWindowCounterMemcachedRepository))
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
	reqCtx := req.Context()
	if r.localCacheRateLimiter.IsAllowed(reqCtx, req) {
		log.Log("passed local checking")
		if r.memcachedRateLimiter.IsAllowed(reqCtx, req) {
			log.Log("passed memcached checking")
			r.next.ServeHTTP(rw, req)
			return
		}
	}
	log.Log("failed local checking")
	rw.WriteHeader(http.StatusTooManyRequests)
	encoder.Encode(map[string]any{"status_code": http.StatusTooManyRequests, "message": "rate limit exceeded, try again later"})
	return
}
