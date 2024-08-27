package fetcher

import (
	"context"
	"github.com/dylanbernhardt/beatradar/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestChromeDPScraperFetchSongs(t *testing.T) {
	// Load the sample HTML
	sampleHTML, err := os.ReadFile(filepath.Join("testdata", "beatstats_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read sample HTML: %v", err)
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(sampleHTML)
	}))
	defer server.Close()

	// Create a ChromeDPScraper with the test server URL
	scraper := NewChromeDPScraper(server.URL)

	// Fetch songs
	ctx := context.Background()
	songs, err := scraper.fetchSongsFromURL(ctx, server.URL)
	if err != nil {
		t.Fatalf("Failed to fetch songs: %v", err)
	}

	// Check the results
	if len(songs) == 0 {
		t.Fatal("No songs were fetched")
	}

	// Check the first song
	expectedFirstSong := models.Song{
		Title:  "Children",
		Artist: "ROBERT MILES, TINLICKER",
		URL:    server.URL + "/track/children/14269418",
	}

	if songs[0].Title != expectedFirstSong.Title {
		t.Errorf("Expected first song title to be %q, got %q", expectedFirstSong.Title, songs[0].Title)
	}
	if songs[0].Artist != expectedFirstSong.Artist {
		t.Errorf("Expected first song artist to be %q, got %q", expectedFirstSong.Artist, songs[0].Artist)
	}
	if songs[0].URL != expectedFirstSong.URL {
		t.Errorf("Expected first song URL to be %q, got %q", expectedFirstSong.URL, songs[0].URL)
	}

	// You can add more checks for other songs if needed
}
