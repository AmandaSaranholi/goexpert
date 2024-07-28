package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/entity"
)

type WeatherRepository interface {
	GetWeatherByLocation(location string) (*entity.Weather, error)
	GetLocationByZipCode(zipCode string) (string, error)
}

type weatherRepository struct {
	viaCepAPI     string
	weatherAPI    string
	weatherAPIKey string
}

func NewWeatherRepository(viaCepAPI, weatherAPI, weatherAPIKey string) WeatherRepository {
	return &weatherRepository{
		viaCepAPI:     viaCepAPI,
		weatherAPI:    weatherAPI,
		weatherAPIKey: weatherAPIKey,
	}
}

func (r *weatherRepository) GetWeatherByLocation(location string) (*entity.Weather, error) {
	url := fmt.Sprintf(r.weatherAPI, r.weatherAPIKey, location)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather data")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	current := result["current"].(map[string]interface{})
	tempC := current["temp_c"].(float64)
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	weather := &entity.Weather{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	return weather, nil
}

func (r *weatherRepository) GetLocationByZipCode(zipCode string) (string, error) {
	url := fmt.Sprintf(r.viaCepAPI, zipCode)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("internal server error")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result["erro"] == "true" {
		return "", fmt.Errorf("can not find zipcode")
	}

	location := fmt.Sprintf("%s,%s", result["localidade"], result["uf"])
	return location, nil
}
