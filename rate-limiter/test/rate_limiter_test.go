package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaSaranholi/goexpert/rate-limiter/config"
	"github.com/AmandaSaranholi/goexpert/rate-limiter/internal/middleware"
	"github.com/AmandaSaranholi/goexpert/rate-limiter/internal/service"
	"github.com/AmandaSaranholi/goexpert/rate-limiter/pkg/ratelimiter"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	cfg := config.Config{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		IPLimit:       5,
		TokenLimit:    5,
		BlockDuration: 60 * time.Second,
		StrategyType:  "redis",
	}

	var repo ratelimiter.Strategy

	switch cfg.StrategyType {
	case "redis":
		repo = ratelimiter.NewRedisStrategy(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	default:
		panic("Unsupported STRATEGY_TYPE")
	}

	rateLimiterService := service.NewRateLimiterService(repo, cfg.IPLimit, cfg.TokenLimit, cfg.BlockDuration)

	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)
	r.Use(middleware.RateLimiterMiddleware(rateLimiterService))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the rate limited !!"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	client := &http.Client{}

	t.Run("Test IP", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			resp, err := client.Do(req)
			assert.NoError(t, err)

			if resp.StatusCode == http.StatusTooManyRequests {
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
				break
			} else {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
		}
	})

	t.Run("Test API_KEY", func(t *testing.T) {
		req.Header.Set("API_KEY", "test-token")

		for i := 0; i < 10; i++ {
			resp, err := client.Do(req)
			assert.NoError(t, err)
			if i < 5 {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}
	})
}
