package traefik_ratelimit

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/bradfitz/gomemcache/memcache"
	slidingWindowCounterRepo "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
	"github.com/ghnexpress/traefik-ratelimit/sliding_window_counter"
	"net/http"
)

type RateLimiter interface {
	IsAllowed(ctx context.Context, req *http.Request) bool
}

type MemcachedConfig struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	MaxRequestInWindow int             `json:"max_request_in_window,omitempty"`
	WindowTime         int             `json:"window_time,omitempty"`
	MemcachedConfig    MemcachedConfig `json:"memcached_config"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type RateLimit struct {
	name        string
	next        http.Handler
	rate        int
	rateLimiter RateLimiter
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	memcacheInstance := memcache.New(config.MemcachedConfig.Address)
	slidingWindowCounterRepository := slidingWindowCounterRepo.NewSlidingWindowCounterRepository(memcacheInstance)
	rateLimiter := sliding_window_counter.NewSlidingWindowCounter(
		slidingWindowCounterRepository,
		sliding_window_counter.SlidingWindowCounterParam{
			MaxRequestInWindow: config.MaxRequestInWindow,
			WindowTime:         config.WindowTime,
		},
	)
	return &RateLimit{
		name:        name,
		next:        next,
		rateLimiter: rateLimiter,
	}, nil
}

func (r *RateLimit) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.rateLimiter.IsAllowed(req.Context(), req) {
		r.next.ServeHTTP(rw, req)
	}
}
