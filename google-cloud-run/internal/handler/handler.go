package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/usecase"
	"github.com/go-chi/chi/v5"
)

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

	weather, err := h.useCase.GetWeatherByZipCode(zipCode)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}
