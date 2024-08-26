package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dylanbernhardt/beatradar/internal/models"
)

func loadTestHTML(t *testing.T) string {
	t.Helper()
	path := filepath.Join("testdata", "beatstats_sample.html")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test HTML file: %v", err)
	}
	return string(content)
}

func TestBeatstatsScraperFetchSongs(t *testing.T) {
	testHTML := loadTestHTML(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewBeatstatsScraper(server.URL)

	songs, err := scraper.FetchSongs(context.Background(), "House", time.Now())
	if err != nil {
		t.Fatalf("FetchSongs failed: %v", err)
	}

	// TODO Add assertions based on the actual content from Beatstats
	if len(songs) == 0 {
		t.Fatalf("Expected to find songs, but got none")
	}

	// TODO Check for a specific song you know should be in the results
	// TODO adjust these based on the actual content of your sample
	expectedSong := models.Song{
		Title:  "Expected Song Title",
		Artist: "Expected Artist",
		// TODO Add other fields
	}

	found := false
	for _, song := range songs {
		if song.Title == expectedSong.Title && song.Artist == expectedSong.Artist {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Did not find the expected song: %v", expectedSong)
	}

	// TODO Add more specific checks based on the actual data
}

func TestBeatstatsScraperFetchSongDetails(t *testing.T) {
	// TODO create a mock server that returns song details HTML
	// TODO Then test the FetchSongDetails method
}
