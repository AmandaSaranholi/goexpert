package http

import (
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *handler.WeatherHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/weather/{zipcode}", handler.GetWeather)
	return r
}
