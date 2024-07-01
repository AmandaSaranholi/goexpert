package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStrategy struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStrategy(addr, password string, db int) *RedisStrategy {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStrategy{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisStrategy) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisStrategy) Expire(key string, duration time.Duration) (bool, error) {
	return r.client.Expire(r.ctx, key, duration).Result()
}

func (r *RedisStrategy) SetBlocked(key string, duration time.Duration) error {
	return r.client.Set(r.ctx, key, "1", duration).Err()
}

func (r *RedisStrategy) IsBlocked(key string) (bool, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
