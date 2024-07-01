package ratelimiter

import "time"

type Strategy interface {
	Increment(key string) (int64, error)
	Expire(key string, duration time.Duration) (bool, error)
	SetBlocked(key string, duration time.Duration) error
	IsBlocked(key string) (bool, error)
}
