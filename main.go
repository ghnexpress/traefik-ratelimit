package traefik_ratelimit

import (
	"context"
	"github.com/ghnexpress/traefik-ratelimit/sliding_window_counter"
	"net/http"

	"github.com/hoisie/redis"
)

type RateLimiter interface {
	IsAllowed(ctx context.Context, req *http.Request) bool
}

var redisClient *redis.Client

// Redis config.
type RedisConfig struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	MaxRequestInWindow int         `json:"max_request_in_window,omitempty"`
	WindowTime         int         `json:"window_time,omitempty"`
	Redis              RedisConfig `json:"redis,omitempty"`
	RedisClient        redis.Client
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

func getRedisClient(config RedisConfig) redis.Client {
	if redisClient == nil {
		redisClient = &redis.Client{
			Addr:        config.Address,
			Db:          0,
			Password:    config.Password,
			MaxPoolSize: 0,
		}
	}

	return *redisClient
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	//redisInstance := getRedisClient(config.Redis)
	rateLimiter := sliding_window_counter.NewSlidingWindowCounter()
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
