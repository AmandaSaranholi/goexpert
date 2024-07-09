package test

import (
	"crypto/rand"
	"encoding/base64"
	"io"
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

func generateRandomToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "token-default"
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

func TestRateLimiter(t *testing.T) {
	cfg := config.Config{
		RedisAddr:     "redis:6379",
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

	client := &http.Client{}

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	t.Run("Test IP", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp, err := client.Do(req)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			if resp.StatusCode == http.StatusTooManyRequests {
				assert.Equal(t, string(body), "you have reached the maximum number of requests or actions allowed within a certain time frame\n")
			} else {
				assert.Equal(t, string(body), "Welcome to the rate limited !!")
			}
		}
	})

	t.Run("Test API_KEY", func(t *testing.T) {
		token := generateRandomToken(10)
		req.Header.Set("API_KEY", token)

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
