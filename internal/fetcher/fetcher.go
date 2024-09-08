package fetcher

import (
	"context"
	"time"

	"github.com/dylanbernhardt/beatradar/internal/models"
)

type SongFetcher interface {
	FetchSongs(ctx context.Context, genre string, date time.Time) ([]models.Song, error)

	FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error)

	FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error)
}
