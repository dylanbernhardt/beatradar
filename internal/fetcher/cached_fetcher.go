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

	return &CachedSongFetcher{
		fetcher: fetcher,
		cache:   redis.NewClient(opts),
		ttl:     ttl,
	}, nil
}

func (c *CachedSongFetcher) FetchSongs(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
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
	songs, err := c.fetcher.FetchSongs(ctx, genre, date)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if cachedData, err := json.Marshal(songs); err == nil {
		c.cache.Set(ctx, cacheKey, cachedData, c.ttl)
	}

	return songs, nil
}

func (c *CachedSongFetcher) FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error) {
	cacheKey := fmt.Sprintf("song_details:%s", songURL)

	// Try to get from cache
	cachedData, err := c.cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var song models.Song
		if err := json.Unmarshal(cachedData, &song); err == nil {
			return &song, nil
		}
	}

	// If not in cache, fetch from the original source
	song, err := c.fetcher.FetchSongDetails(ctx, songURL)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if cachedData, err := json.Marshal(song); err == nil {
		c.cache.Set(ctx, cacheKey, cachedData, c.ttl)
	}

	return song, nil
}

func (c *CachedSongFetcher) FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	cacheKey := fmt.Sprintf("songs_with_details:%s:%s", genre, date.Format("2006-01-02"))

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
