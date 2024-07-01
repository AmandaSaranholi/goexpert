package service

import (
	"time"

	"github.com/AmandaSaranholi/goexpert/rate-limiter/pkg/ratelimiter"
)

type RateLimiter interface {
	AllowRequest(ip string, token string) (bool, string)
}

type RateLimiterService struct {
	repo          ratelimiter.Strategy
	ipLimit       int
	tokenLimit    int
	blockDuration time.Duration
}

func NewRateLimiterService(repo ratelimiter.Strategy, ipLimit, tokenLimit int, blockDuration time.Duration) *RateLimiterService {
	return &RateLimiterService{
		repo:          repo,
		ipLimit:       ipLimit,
		tokenLimit:    tokenLimit,
		blockDuration: blockDuration,
	}
}

func (rl *RateLimiterService) AllowRequest(ip, token string) (bool, string) {
	key := "ip:" + ip
	limit := rl.ipLimit

	if token != "" {
		key = "token:" + token
		limit = rl.tokenLimit
	}

	blocked, err := rl.repo.IsBlocked("blocked:" + key)
	if err != nil {
		return false, "Internal Server Error"
	}
	if blocked {
		return false, "you have reached the maximum number of requests or actions allowed within a certain time frame"
	}

	count, err := rl.repo.Increment(key)
	if err != nil {
		return false, "Internal Server Error"
	}

	if count > int64(limit) {
		rl.repo.SetBlocked("blocked:"+key, rl.blockDuration)
		return false, "you have reached the maximum number of requests or actions allowed within a certain time frame"
	}

	rl.repo.Expire(key, time.Second)
	return true, ""
}
