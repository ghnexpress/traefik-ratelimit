package rate_limiter

import (
	"context"
	"net/http"
)

type RateLimiter interface {
	IsAllowed(ctx context.Context, req *http.Request, rw http.ResponseWriter) bool
}
