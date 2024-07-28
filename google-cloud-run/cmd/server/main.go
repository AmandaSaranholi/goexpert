package main

import (
	"log"
	"net/http"

	httpRouter "github.com/AmandaSaranholi/goexpert/google-cloud-run/pkg/http"

	"github.com/AmandaSaranholi/goexpert/google-cloud-run/config"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/handler"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/repository"
	"github.com/AmandaSaranholi/goexpert/google-cloud-run/internal/usecase"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	repo := repository.NewWeatherRepository(cfg.ViaCepAPI, cfg.WeatherAPI, cfg.WeatherAPIKey)
	useCase := usecase.NewWeatherUseCase(repo)
	handler := handler.NewWeatherHandler(useCase)
	router := httpRouter.NewRouter(handler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
