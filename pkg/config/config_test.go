package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "9090")
	os.Setenv("BEATSTATS_URL", "https://test.beatstats.com")
	os.Setenv("REDIS_URL", "redis://testhost:6379")
	os.Setenv("CACHE_TTL", "3600")
	os.Setenv("SCRAPER_TIMEOUT", "60")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Port)
	}
	if cfg.BeatstatsURL != "https://test.beatstats.com" {
		t.Errorf("Expected BeatstatsURL to be 'https://test.beatstats.com', got '%s'", cfg.BeatstatsURL)
	}
	if cfg.RedisURL != "redis://testhost:6379" {
		t.Errorf("Expected RedisURL to be 'redis://testhost:6379', got '%s'", cfg.RedisURL)
	}
	if cfg.CacheTTL != time.Hour {
		t.Errorf("Expected CacheTTL to be 1 hour, got %v", cfg.CacheTTL)
	}
	if cfg.ScraperTimeout != time.Minute {
		t.Errorf("Expected ScraperTimeout to be 1 minute, got %v", cfg.ScraperTimeout)
	}
}

func TestGetGenreName(t *testing.T) {
	tests := []struct {
		id       string
		expected string
	}{
		{"7", "House"},
		{"8", "Techno"},
		{"999", "Unknown Genre"},
	}

	for _, test := range tests {
		result := GetGenreName(test.id)
		if result != test.expected {
			t.Errorf("For id '%s', expected '%s', got '%s'", test.id, test.expected, result)
		}
	}
}

func TestGetGenreID(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		found    bool
	}{
		{"House", "7", true},
		{"Techno", "8", true},
		{"Unknown", "", false},
	}

	for _, test := range tests {
		result, ok := GetGenreID(test.name)
		if ok != test.found {
			t.Errorf("For genre '%s', expected found to be %v, got %v", test.name, test.found, ok)
		}
		if result != test.expected {
			t.Errorf("For genre '%s', expected id '%s', got '%s'", test.name, test.expected, result)
		}
	}
}
