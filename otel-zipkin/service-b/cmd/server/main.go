package main

import (
	"context"
	"log"
	"net/http"
	"os"

	httpRouter "github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/pkg/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/config"
	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/handler"
	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/repository"
	"github.com/AmandaSaranholi/goexpert/otel-zipkin/service-b/internal/usecase"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	collectorURL := os.Getenv("OTEL_EXPORTER_URL")
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(collectorURL), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("orchestrator-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func main() {
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		panic(err)
	}

	tp, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	repo := repository.NewWeatherRepository(cfg.ViaCepAPI, cfg.WeatherAPI, cfg.WeatherAPIKey)
	useCase := usecase.NewWeatherUseCase(repo)
	handler := handler.NewWeatherHandler(useCase)
	router := httpRouter.NewRouter(handler)

	log.Println("Starting server on :8181")
	http.ListenAndServe(":8181", router)
}
