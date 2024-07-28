package usecase

import (
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/entity"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/repository"
)

type WeatherUseCase interface {
	GetWeatherByZipCode(zipCode string) (*entity.Weather, error)
}

type weatherUseCase struct {
	repo repository.WeatherRepository
}

func NewWeatherUseCase(repo repository.WeatherRepository) WeatherUseCase {
	return &weatherUseCase{repo: repo}
}

func (u *weatherUseCase) GetWeatherByZipCode(zipCode string) (*entity.Weather, error) {
	location, err := u.repo.GetLocationByZipCode(zipCode)
	if err != nil {
		return nil, err
	}

	weather, err := u.repo.GetWeatherByLocation(location)
	if err != nil {
		return nil, err
	}

	return weather, nil
}
