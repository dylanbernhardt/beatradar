package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dylanbernhardt/beatradar/internal/models"
	"github.com/dylanbernhardt/beatradar/pkg/config"
)

type ChromeDPScraper struct {
	baseURL string
}

func NewChromeDPScraper(baseURL string) *ChromeDPScraper {
	return &ChromeDPScraper{baseURL: baseURL}
}

func (c *ChromeDPScraper) FetchSongs(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	genreID, ok := config.GetGenreID(genre)
	if !ok {
		return nil, fmt.Errorf("unknown genre: %s", genre)
	}

	url := fmt.Sprintf("%s/tracks/home/list?genre=%s&period=2&datefilter=%s",
		c.baseURL, genreID, date.Format("2006-01-02"))

	return c.fetchSongsFromURL(ctx, url)
}

func (c *ChromeDPScraper) fetchSongsFromURL(ctx context.Context, url string) ([]models.Song, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var songs []models.Song

	err := chromedp.Run(ctx,
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

	if err != nil {
		return nil, fmt.Errorf("error scraping with ChromeDP: %w", err)
	}

	return songs, nil
}

func (c *ChromeDPScraper) FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error) {
	// Implement if needed
	return nil, fmt.Errorf("FetchSongDetails not implemented for ChromeDPScraper")
}

func (c *ChromeDPScraper) FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	return c.FetchSongs(ctx, genre, date)
}
