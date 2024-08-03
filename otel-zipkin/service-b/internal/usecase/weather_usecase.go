package usecase

import (
	"context"

	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/entity"
	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/repository"
)

type WeatherUseCase interface {
	GetWeatherByZipCode(zipCode string, ctx context.Context) (*entity.Weather, error)
}

type weatherUseCase struct {
	repo repository.WeatherRepository
}

func NewWeatherUseCase(repo repository.WeatherRepository) WeatherUseCase {
	return &weatherUseCase{repo: repo}
}

func (u *weatherUseCase) GetWeatherByZipCode(zipCode string, ctx context.Context) (*entity.Weather, error) {
	location, err := u.repo.GetLocationByZipCode(zipCode, ctx)
	if err != nil {
		return nil, err
	}

	weather, err := u.repo.GetWeatherByLocation(location, ctx)
	if err != nil {
		return nil, err
	}

	return weather, nil
}
