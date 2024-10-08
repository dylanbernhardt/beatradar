package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dylanbernhardt/beatradar/internal/cache"
	"github.com/dylanbernhardt/beatradar/internal/models"
)

type CachedFetcher struct {
	fetcher SongFetcher
	cache   *cache.RedisClient
	ttl     time.Duration
}

func NewCachedFetcher(fetcher SongFetcher, cache *cache.RedisClient, ttl time.Duration) *CachedFetcher {
	return &CachedFetcher{
		fetcher: fetcher,
		cache:   cache,
		ttl:     ttl,
	}
}

func (cf *CachedFetcher) FetchSongs(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	cacheKey := fmt.Sprintf("songs:%s:%s", genre, date.Format("2006-01-02"))

	cachedData, err := cf.cache.Get(ctx, cacheKey)
	if err == nil {
		var songs []models.Song
		if err := json.Unmarshal([]byte(cachedData), &songs); err == nil {
			return songs, nil
		}
	}

	songs, err := cf.fetcher.FetchSongs(ctx, genre, date)
	if err != nil {
		return nil, err
	}

	if cachedData, err := json.Marshal(songs); err == nil {
		err := cf.cache.Set(ctx, cacheKey, cachedData, cf.ttl)
		if err != nil {
			return nil, err
		}
	}

	return songs, nil
}

// TODO Implement other methods of the SongFetcher interface...
