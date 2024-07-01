package main

import (
	"fmt"
	"net/http"

	"github.com/AmandaSaranholi/goexpert/rate-limiter/config"
	"github.com/AmandaSaranholi/goexpert/rate-limiter/internal/middleware"
	"github.com/AmandaSaranholi/goexpert/rate-limiter/internal/service"
	"github.com/AmandaSaranholi/goexpert/rate-limiter/pkg/ratelimiter"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig(nil)

	var repo ratelimiter.Strategy

	switch cfg.StrategyType {
	case "redis":
		repo = ratelimiter.NewRedisStrategy(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	default:
		panic("Unsupported STRATEGY_TYPE")
	}

	rateLimiterService := service.NewRateLimiterService(repo, cfg.IPLimit, cfg.TokenLimit, cfg.BlockDuration)

	r := chi.NewRouter()
	r.Use(middleware.RateLimiterMiddleware(rateLimiterService))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the rate limited !!"))
	})

	fmt.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", r)
}
