package sliding_window_counter

import (
	"context"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter"
	slidingWindowCounterRepo "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
	"github.com/ghnexpress/traefik-ratelimit/utils"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	seperatedPart = 60 // should be less than 100
)

func getIp(req *http.Request) string {
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			return strings.TrimSpace(ip[:i])
		}
		return ip
	}
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ra
}

type SlidingWindowCounterParam struct {
	MaxRequestInWindow int
	WindowTime         int
}

type slidingWindowCounter struct {
	repo   slidingWindowCounterRepo.Repository
	params SlidingWindowCounterParam
}

func NewSlidingWindowCounter(repo slidingWindowCounterRepo.Repository, params SlidingWindowCounterParam) rate_limiter.RateLimiter {
	return &slidingWindowCounter{repo: repo, params: params}
}

func (s *slidingWindowCounter) getCurrentPart() int {
	var currPart int
	now := time.Now().Unix()
	periodOfPart := float64(s.params.WindowTime) / float64(seperatedPart)
	currPart = int(math.Round(float64(now)/periodOfPart) * periodOfPart)
	return currPart
}

func (s *slidingWindowCounter) isIPExist(ctx context.Context, ip string) bool {
	_, err := s.repo.GetRequestCountByIP(ctx, ip)
	if err != nil {
		return false
	}
	return true
}

func (s *slidingWindowCounter) increaseAndGetTotalRequestInWindow(ctx context.Context, ip string, part int) (cumulativeReq int, err error) {
	errChan := make(chan error, 2)
	defer close(errChan)
	w := sync.WaitGroup{}
	go func() {
		w.Add(1)
		defer w.Done()
		if err := s.repo.RemoveExpiredWindowSlice(ctx, ip, part, s.params.WindowTime); err != nil {
			utils.ShowErrorLogs(err)
			errChan <- err
		}
	}()

	go func() {
		w.Add(1)
		defer w.Done()
		if err := s.repo.IncreaseCurrentWindowSlice(ctx, ip, part); err != nil {
			utils.ShowErrorLogs(err)
			errChan <- err
		}
	}()

	w.Wait()
	if err = <-errChan; err != nil {
		return 0, err
	}

	cumulativeReq, err = s.repo.GetAllRequestCountCurrentWindow(ctx, ip)
	if err != nil {
		return 0, err
	}

	return cumulativeReq, err
}

func (s *slidingWindowCounter) IsAllowed(ctx context.Context, req *http.Request) bool {
	reqIP := getIp(req)
	if !s.isIPExist(ctx, reqIP) {
		if err := s.repo.AddNewIP(ctx, reqIP); err != nil {
			utils.ShowErrorLogs(err)
			return false
		}
	}

	currPart := s.getCurrentPart()
	cumulativeReq, err := s.increaseAndGetTotalRequestInWindow(ctx, reqIP, currPart)
	if err != nil {
		utils.ShowErrorLogs(err)
		return false
	}
	if cumulativeReq > s.params.MaxRequestInWindow {
		return false
	}
	return true
}
