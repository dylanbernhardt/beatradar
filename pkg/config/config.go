package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	BeatstatsURL   string
	RedisURL       string
	CacheTTL       time.Duration
	ScraperTimeout time.Duration
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	beatstatsURL := os.Getenv("BEATSTATS_URL")
	if beatstatsURL == "" {
		beatstatsURL = "https://api.beatradar.app" // Updated URL
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379" // Default Redis URL
	}

	cacheTTL, err := getEnvDuration("CACHE_TTL", 24*time.Hour) // Default 24 hours
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL: %w", err)
	}

	scraperTimeout, err := getEnvDuration("SCRAPER_TIMEOUT", 30*time.Second) // Default 30 seconds
	if err != nil {
		return nil, fmt.Errorf("invalid SCRAPER_TIMEOUT: %w", err)
	}

	return &Config{
		Port:           port,
		BeatstatsURL:   beatstatsURL,
		RedisURL:       redisURL,
		CacheTTL:       cacheTTL,
		ScraperTimeout: scraperTimeout,
	}, nil
}

func getEnvDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration: %s", value)
	}

	return time.Duration(intValue) * time.Second, nil
}
