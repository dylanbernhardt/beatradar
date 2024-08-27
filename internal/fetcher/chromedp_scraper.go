package fetcher

import (
	"context"
	"fmt"
	"log"
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
	return c.fetchSongsFromURL(ctx, c.constructURL(genre, date))
}

func (c *ChromeDPScraper) constructURL(genre string, date time.Time) string {
	genreID, ok := config.GetGenreID(genre)
	if !ok {
		log.Printf("Unknown genre: %s, using as is", genre)
		genreID = genre
	}
	return fmt.Sprintf("%s/tracks/home/list?genre=%s&period=2&datefilter=%s",
		c.baseURL, genreID, date.Format("2006-01-02"))
}

func (c *ChromeDPScraper) fetchSongsFromURL(ctx context.Context, url string) ([]models.Song, error) {
	// Create a new chrome instance
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
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
