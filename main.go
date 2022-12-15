package traefik_ratelimit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hoisie/redis"
)

var redisClient *redis.Client

// Redis config.
type RedisConfig struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	Rate        int         `json:"rate,omitempty"`
	Redis       RedisConfig `json:"redis,omitempty"`
	RedisClient redis.Client
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type RateLimit struct {
	name  string
	next  http.Handler
	rate  int
	redis redis.Client
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
	return &RateLimit{
		name:  name,
		next:  next,
		rate:  config.Rate,
		redis: getRedisClient(config.Redis),
	}, nil
}

func (r *RateLimit) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := req.RemoteAddr
	val, _ := r.redis.Incr(key)
	req.Header.Set("Count", fmt.Sprintf("%d", val))
	req.Header.Set("Rate", fmt.Sprintf("%d", r.rate))
	req.Header.Set("Version", "v1.0.0")
	r.next.ServeHTTP(rw, req)
}
