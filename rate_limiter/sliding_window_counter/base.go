package sliding_window_counter

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter"
	slidingWindowCounterRepo "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
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
	repo           slidingWindowCounterRepo.Repository
	errorPublisher log.ErrorPublisher
	params         SlidingWindowCounterParam
}

func NewSlidingWindowCounter(repo slidingWindowCounterRepo.Repository, errorPublisher log.ErrorPublisher, params SlidingWindowCounterParam) rate_limiter.RateLimiter {
	return &slidingWindowCounter{repo: repo, errorPublisher: errorPublisher, params: params}
}

func (s *slidingWindowCounter) getFormattedError(ctx context.Context, err error, extra string) error {
	env := "dev"
	if value, ok := ctx.Value("env").(string); ok {
		env = value
	}
	requestID := ""
	if value, ok := ctx.Value("requestID").(string); ok {
		requestID = value
	}
	return fmt.Errorf("[%s][rate-limiter-middleware-plugin] %s \nRequestID: %s\n%s", env, extra, requestID, err.Error())
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
		log.Log("err ", err)
		return false
	}
	return true
}

func (s *slidingWindowCounter) increaseAndGetTotalRequestInWindow(ctx context.Context, ip string, part int) (cumulativeReq int, err error) {
	defer func() {
		if err := recover(); err != nil {
			errRes := s.getFormattedError(ctx, err.(error), "")
			go s.errorPublisher.SendError(errRes)
			log.Log(err)
		}
	}()
	errChan := make(chan error, 2)
	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		defer w.Done()
		if err := s.repo.RemoveExpiredWindowSlice(ctx, ip, part, s.params.WindowTime); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer w.Done()
		if err := s.repo.IncreaseCurrentWindowSlice(ctx, ip, part); err != nil {
			errChan <- err
		}
	}()

	w.Wait()
	close(errChan)
	// 2 lines below are important, if change this to below code can lead to panic
	// if err = <-errChan; err != nil {
	//		return 0, err
	//	}
	if err := <-errChan; err != nil {
		return 0, err
	}
	cumulativeReq, err = s.repo.GetAllRequestCountCurrentWindow(ctx, ip)
	if err != nil {
		return 0, err
	}
	return cumulativeReq, err
}

func (s *slidingWindowCounter) IsAllowed(ctx context.Context, req *http.Request) bool {
	defer func() {
		if err := recover(); err != nil {
			errRes := s.getFormattedError(ctx, err.(error), "")
			go s.errorPublisher.SendError(errRes)
			log.Log(err)
		}
	}()
	reqIP := getIp(req)
	if !s.isIPExist(ctx, reqIP) {
		if err := s.repo.AddNewIP(ctx, reqIP); err != nil {
			err = s.getFormattedError(ctx, err, "add new ip error")
			go s.errorPublisher.SendError(err)
			return false
		}
	}

	currPart := s.getCurrentPart()

	cumulativeReq, err := s.increaseAndGetTotalRequestInWindow(ctx, reqIP, currPart)
	if err != nil {
		err = s.getFormattedError(ctx, err, "increase and get total request in window")
		go s.errorPublisher.SendError(err)
		return false
	}

	if cumulativeReq > s.params.MaxRequestInWindow {
		return false
	}
	return true
}
