package sliding_window_counter

import (
	"context"
	"fmt"
	"github.com/ghnexpress/traefik-ratelimit/utils"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/rate_limiter"
	slidingWindowCounterRepo "github.com/ghnexpress/traefik-ratelimit/repo/sliding_window_counter"
)

const (
	seperatedPart = 60 // should be less than 100
)

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
		return false
	}
	return true
}

func (s *slidingWindowCounter) increaseAndGetTotalRequestInWindow(ctx context.Context, ip string, part int) (cumulativeReq int, err error) {
	defer func() {
		if err := recover(); err != nil {
			errRes := s.getFormattedError(ctx, err.(error), "")

			go s.errorPublisher.SendError(errRes)
		}
	}()
	errChan := make(chan error, 2)
	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		defer w.Done()
		if err := s.repo.RemoveExpiredWindowSlice(ctx, ip, part, s.params.WindowTime); err != nil {
			errChan <- fmt.Errorf("remove expired window slice err %v - current part %d", err, part)
		}
	}()

	go func() {
		defer w.Done()
		if err := s.repo.IncreaseCurrentWindowSlice(ctx, ip, part); err != nil {
			errChan <- fmt.Errorf("increase current window slice err %v - current part %d", err, part)
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
		return 0, fmt.Errorf("get all request count of current window err %v", err)
	}
	return cumulativeReq, err
}

func (s *slidingWindowCounter) IsAllowed(ctx context.Context, req *http.Request, rw http.ResponseWriter) bool {
	defer func() {
		if err := recover(); err != nil {
			errRes := s.getFormattedError(ctx, err.(error), "")
			go s.errorPublisher.SendError(errRes)
			log.Log(err)
		}
	}()
	reqIP := utils.GetIp(req)
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
		err = s.getFormattedError(ctx, err, fmt.Sprintf("%s increase and get total request in window", reqIP))
		go s.errorPublisher.SendError(err)
		return false
	}
	rw.Header().Add("Num-request", strconv.Itoa(cumulativeReq))
	if cumulativeReq > s.params.MaxRequestInWindow {
		return false
	}

	return true
}
