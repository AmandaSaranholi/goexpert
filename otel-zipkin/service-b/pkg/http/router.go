package http

import (
	"net/http"

	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func NewRouter(handler *handler.WeatherHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "chi-server", otelhttp.WithPropagators(otel.GetTextMapPropagator()))
	})
	r.Get("/weather/{zipcode}", handler.GetWeather)
	return r
}
