package traefik_ratelimit

import (
	"context"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter/sliding_window_counter"
	slidingWindowCounterRepo "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
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
	name        string
	next        http.Handler
	rate        int
	rateLimiter rate_limiter.RateLimiter
	config      *Config
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log.Log(fmt.Sprintf("config %v", config))
	memcachedInstance := memcache.New(config.MemcachedConfig.Address)

	slidingWindowCounterRepository := slidingWindowCounterRepo.NewSlidingWindowCounterRepository(memcachedInstance)
	rateLimiter := sliding_window_counter.NewSlidingWindowCounter(
		slidingWindowCounterRepository,
		sliding_window_counter.SlidingWindowCounterParam{
			MaxRequestInWindow: config.MaxRequestInWindow,
			WindowTime:         config.WindowTime,
		},
	)
	log.Log(fmt.Sprintf("%v", slidingWindowCounterRepository))
	return &RateLimit{
		name:        name,
		next:        next,
		rateLimiter: rateLimiter,
		config:      config,
	}, nil
}

func (r *RateLimit) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.rateLimiter.IsAllowed(req.Context(), req) {
		r.next.ServeHTTP(rw, req)
	} else {
		http.Error(rw, fmt.Sprintf("rate limit exceeded, try again in %d seconds", r.config.WindowTime), http.StatusTooManyRequests)
	}
}
