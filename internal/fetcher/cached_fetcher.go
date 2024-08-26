package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dylanbernhardt/beatradar/internal/models"
	"github.com/go-redis/redis/v8"
)

type CachedSongFetcher struct {
	fetcher SongFetcher
	cache   *redis.Client
	ttl     time.Duration
}

func NewCachedSongFetcher(fetcher SongFetcher, redisURL string, ttl time.Duration) (*CachedSongFetcher, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	return &CachedSongFetcher{
		fetcher: fetcher,
		cache:   client,
		ttl:     ttl,
	}, nil
}

func (c *CachedSongFetcher) FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	cacheKey := fmt.Sprintf("songs:%s:%s", genre, date.Format("2006-01-02"))

	// Try to get from cache
	cachedData, err := c.cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var songs []models.Song
		if err := json.Unmarshal(cachedData, &songs); err == nil {
			return songs, nil
		}
	}

	// If not in cache, fetch from the original source
	songs, err := c.fetcher.FetchSongsWithDetails(ctx, genre, date)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if cachedData, err := json.Marshal(songs); err == nil {
		c.cache.Set(ctx, cacheKey, cachedData, c.ttl)
	}

	return songs, nil
}

// Close closes the Redis connection
func (c *CachedSongFetcher) Close() error {
	return c.cache.Close()
}
