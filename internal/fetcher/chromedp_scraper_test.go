package fetcher

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/dylanbernhardt/beatradar/internal/models"
)

func TestChromeDPScraperFetchSongs(t *testing.T) {
	// Get the absolute path to the sample HTML file
	absPath, err := filepath.Abs(filepath.Join("testdata", "beatstats_sample.html"))
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a file URL
	fileURL := "file://" + absPath

	// Create a ChromeDPScraper with an empty base URL (not used in this test)
	scraper := NewChromeDPScraper("")

	// Fetch songs
	ctx := context.Background()
	songs, err := scraper.fetchSongsFromURL(ctx, fileURL)
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
	}

	if songs[0].Title != expectedFirstSong.Title {
		t.Errorf("Expected first song title to be %q, got %q", expectedFirstSong.Title, songs[0].Title)
	}
	if songs[0].Artist != expectedFirstSong.Artist {
		t.Errorf("Expected first song artist to be %q, got %q", expectedFirstSong.Artist, songs[0].Artist)
	}

	// You can add more checks for other songs if needed
}
