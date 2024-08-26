package fetcher

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dylanbernhardt/beatradar/internal/models"
	"github.com/dylanbernhardt/beatradar/pkg/config"
)

type SongFetcher interface {
	FetchSongs(ctx context.Context, genreID string, date time.Time) ([]models.Song, error)
	FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error)
	FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error)
}

type BeatstatsScraper struct {
	baseURL    string
	httpClient *http.Client
}

func NewBeatstatsScraper(baseURL string) *BeatstatsScraper {
	return &BeatstatsScraper{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (b *BeatstatsScraper) FetchSongs(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	genreID, ok := config.GetGenreID(genre)
	if !ok {
		return nil, fmt.Errorf("unknown genre: %s", genre)
	}

	url := fmt.Sprintf("%s/tracks/home/list?genre=%s&period=2&datefilter=%s",
		b.baseURL, genreID, date.Format("2006-01-02"))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	var songs []models.Song

	doc.Find(".track-row").Each(func(i int, s *goquery.Selection) {
		songURL, _ := s.Find(".track-title a").Attr("href")
		title := s.Find(".track-title a").Text()
		artist := s.Find(".track-artists").Text()

		song := models.Song{
			Title:  title,
			Artist: artist,
			URL:    b.baseURL + songURL,
			Genre:  config.GenreMap[genreID],
		}

		songs = append(songs, song)
	})

	return songs, nil
}

func (b *BeatstatsScraper) FetchSongDetails(ctx context.Context, songURL string) (*models.Song, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", songURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	song := &models.Song{}

	song.Title = doc.Find(".track-title").Text()
	song.Artist = doc.Find(".track-artists").Text()
	song.Genre = doc.Find(".track-genre").Text()

	// TODO add other details like BPM, Key, Length, etc.
	// TODO adjust these selectors based on the actual HTML structure

	return song, nil
}

func (b *BeatstatsScraper) FetchSongsWithDetails(ctx context.Context, genre string, date time.Time) ([]models.Song, error) {
	songs, err := b.FetchSongs(ctx, genre, date)
	if err != nil {
		return nil, err
	}

	var songsWithDetails []models.Song

	for _, song := range songs {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

		details, err := b.FetchSongDetails(ctx, song.URL)
		if err != nil {
			fmt.Printf("Error fetching details for %s: %v\n", song.Title, err)
			continue
		}

		song.BPM = details.BPM
		song.Key = details.Key
		song.Length = details.Length
		song.ReleaseDate = details.ReleaseDate

		songsWithDetails = append(songsWithDetails, song)
	}

	return songsWithDetails, nil
}
