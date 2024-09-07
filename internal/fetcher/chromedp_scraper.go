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

	if err := ctx.Err(); err != nil {
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
				err = fmt.Errorf("panic occurred during ChromeDP run: %v", r)
			}
		}()

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
		return
	}()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout while scraping with ChromeDP: %w", err)
		}
		if strings.Contains(err.Error(), "net::ERR_") {
			return nil, fmt.Errorf("network error while scraping with ChromeDP: %w", err)
		}
		return nil, fmt.Errorf("error scraping with ChromeDP: %w", err)
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("no songs found at URL: %s", url)
	}

	return songs, nil
}

func (c *ChromeDPScraper) FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error) {
	// TODO Implement if needed
	return nil, fmt.Errorf("FetchSongDetails not implemented for ChromeDPScraper")
}

func (c *ChromeDPScraper) FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	return c.FetchSongs(ctx, genre, date)
}
