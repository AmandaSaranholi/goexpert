package middleware

import (
	"net/http"

	"github.com/AmandaSaranholi/goexpert/rate-limiter/internal/service"
)

func RateLimiterMiddleware(rl *service.RateLimiterService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			token := r.Header.Get("API_KEY")

			allowed, message := rl.AllowRequest(ip, token)
			if !allowed {
				http.Error(w, message, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
