package fetcher

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dylanbernhardt/beatradar/internal/models"
	"github.com/dylanbernhardt/beatradar/pkg/config"
	"go.uber.org/zap"
)

type ChromeDPScraper struct {
	baseURL string
	logger  *zap.Logger
}

func NewChromeDPScraper(baseURL string, logger *zap.Logger) *ChromeDPScraper {
	return &ChromeDPScraper{baseURL: baseURL, logger: logger}
}

func (c *ChromeDPScraper) FetchSongs(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	logger := c.logger.With(
		zap.String("genre", genre),
		zap.Time("date", date),
	)

	genreID, ok := config.GetGenreID(genre)
	if !ok {
		logger.Error("Unknown genre")
		return nil, fmt.Errorf("unknown genre: %s", genre)
	}

	url := fmt.Sprintf("%s/tracks/home/list?genre=%s&period=2&datefilter=%s",
		c.baseURL, genreID, date.Format("2006-01-02"))

	logger.Info("Fetching songs", zap.String("url", url))
	return c.fetchSongsFromURL(ctx, url)
}

func (c *ChromeDPScraper) fetchSongsFromURL(ctx context.Context, url string) ([]models.Song, error) {
	logger := c.logger.With(zap.String("url", url))

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	if err := ctx.Err(); err != nil {
		logger.Error("Parent context is already done", zap.Error(err))
		return nil, fmt.Errorf("parent context is already done: %w", err)
	}

	ctx, cancel = context.WithTimeout(allocCtx, 30*time.Second)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var songs []models.Song

	err := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic occurred during ChromeDP run", zap.Any("panic", r))
				err = fmt.Errorf("panic occurred during ChromeDP run: %v", r)
			}
		}()

		startTime := time.Now()
		err = chromedp.Run(browserCtx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`div[id^='top10artistchart-full']`, chromedp.ByQuery),
			chromedp.Evaluate(`
                Array.from(document.querySelectorAll("div[id^='top10artistchart-full']")).map(el => {
                    return {
                        title: el.querySelector("div[id^='top10trackchart-title']").textContent.trim(),
                        artist: el.querySelector("div[id^='top10trackchart-artistname']").textContent.trim(),
                        url: el.querySelector("a").href
                    }
                })
            `, &songs),
		)
		duration := time.Since(startTime)
		logger.Info("ChromeDP run completed", zap.Duration("duration", duration))
		return
	}()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("Timeout while scraping with ChromeDP", zap.Error(err))
			return nil, fmt.Errorf("timeout while scraping with ChromeDP: %w", err)
		}
		if strings.Contains(err.Error(), "net::ERR_") {
			logger.Error("Network error while scraping with ChromeDP", zap.Error(err))
			return nil, fmt.Errorf("network error while scraping with ChromeDP: %w", err)
		}
		logger.Error("Error scraping with ChromeDP", zap.Error(err))
		return nil, fmt.Errorf("error scraping with ChromeDP: %w", err)
	}

	if len(songs) == 0 {
		logger.Warn("No songs found", zap.String("url", url))
		return nil, fmt.Errorf("no songs found at URL: %s", url)
	}

	logger.Info("Successfully fetched songs", zap.Int("count", len(songs)))
	return songs, nil
}

func (c *ChromeDPScraper) FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error) {
	c.logger.Warn("FetchSongDetails not implemented for ChromeDPScraper", zap.String("songURL", songURL))
	return nil, fmt.Errorf("FetchSongDetails not implemented for ChromeDPScraper")
}

func (c *ChromeDPScraper) FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	return c.FetchSongs(ctx, genre, date)
}
