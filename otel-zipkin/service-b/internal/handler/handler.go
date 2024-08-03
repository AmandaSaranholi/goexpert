package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/usecase"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("orchestrator-service")

type WeatherHandler struct {
	useCase usecase.WeatherUseCase
}

func NewWeatherHandler(useCase usecase.WeatherUseCase) *WeatherHandler {
	return &WeatherHandler{useCase: useCase}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	zipCode := chi.URLParam(r, "zipcode")
	if len(zipCode) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ctx, span := tracer.Start(r.Context(), "get-weather")
	defer span.End()

	weather, err := h.useCase.GetWeatherByZipCode(zipCode, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}
