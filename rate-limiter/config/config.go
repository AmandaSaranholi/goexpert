package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	IPLimit       int
	TokenLimit    int
	BlockDuration time.Duration
	StrategyType  string
}

func LoadConfig(path *string) *Config {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defaultPath := ".env"
	if path == nil {
		path = &defaultPath
	}

	err = godotenv.Load(filepath.Join(pwd, *path))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value")
	}

	ipLimit, err := strconv.Atoi(os.Getenv("LIMITER_IP_LIMIT"))
	if err != nil {
		log.Fatalf("Invalid LIMITER_IP_LIMIT value")
	}

	tokenLimit, err := strconv.Atoi(os.Getenv("LIMITER_TOKEN_LIMIT"))
	if err != nil {
		log.Fatalf("Invalid LIMITER_TOKEN_LIMIT value")
	}

	blockDuration, err := strconv.Atoi(os.Getenv("LIMITER_BLOCK_DURATION"))
	if err != nil {
		log.Fatalf("Invalid LIMITER_BLOCK_DURATION value")
	}

	return &Config{
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		IPLimit:       ipLimit,
		TokenLimit:    tokenLimit,
		BlockDuration: time.Duration(blockDuration) * time.Second,
		StrategyType:  os.Getenv("STRATEGY_TYPE"),
	}
}
