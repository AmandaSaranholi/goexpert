package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AmandaSaranholi/goexpert/google-cloud-run/config"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/handler"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/repository"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/usecase"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather_Success(t *testing.T) {
	r := chi.NewRouter()

	cfg, _ := config.LoadConfig(".")
	repo := repository.NewWeatherRepository(cfg.ViaCepAPI, cfg.WeatherAPI, cfg.WeatherAPIKey)
	useCase := usecase.NewWeatherUseCase(repo)
	handler := handler.NewWeatherHandler(useCase)

	r.Get("/weather/{zipcode}", handler.GetWeather)

	req := httptest.NewRequest(http.MethodGet, "/weather/26572070", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var body map[string]float64
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)

	assert.Contains(t, body, "temp_C")
	assert.Contains(t, body, "temp_F")
	assert.Contains(t, body, "temp_K")
}

func TestGetWeather_InvalidZipcode(t *testing.T) {
	r := chi.NewRouter()

	cfg, _ := config.LoadConfig(".")
	repo := repository.NewWeatherRepository(cfg.ViaCepAPI, cfg.WeatherAPI, cfg.WeatherAPIKey)
	useCase := usecase.NewWeatherUseCase(repo)
	handler := handler.NewWeatherHandler(useCase)

	r.Get("/weather/{zipcode}", handler.GetWeather)

	req := httptest.NewRequest(http.MethodGet, "/weather/123", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	assert.Equal(t, "invalid zipcode\n", resp.Body.String())
}

func TestGetWeather_ZipcodeNotFound(t *testing.T) {
	r := chi.NewRouter()

	cfg, _ := config.LoadConfig(".")
	repo := repository.NewWeatherRepository(cfg.ViaCepAPI, cfg.WeatherAPI, cfg.WeatherAPIKey)
	useCase := usecase.NewWeatherUseCase(repo)
	handler := handler.NewWeatherHandler(useCase)

	r.Get("/weather/{zipcode}", handler.GetWeather)

	req := httptest.NewRequest(http.MethodGet, "/weather/11111111", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, "can not find zipcode\n", resp.Body.String())
}

func TestGetWeather_FetchWeatherFailure(t *testing.T) {
	r := chi.NewRouter()

	cfg, _ := config.LoadConfig(".")
	repo := repository.NewWeatherRepository(cfg.ViaCepAPI, cfg.WeatherAPI, cfg.WeatherAPIKey)
	useCase := usecase.NewWeatherUseCase(repo)
	handler := handler.NewWeatherHandler(useCase)

	r.Get("/weather/{zipcode}", handler.GetWeather)

	req := httptest.NewRequest(http.MethodGet, "/weather/01001000", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	assert.Equal(t, "internal server error\n", resp.Body.String())

}
