package fetcher

import (
	"context"
	"github.com/dylanbernhardt/beatradar/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
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
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewBeatstatsScraper(server.URL)
	scraper.beatportURL = "https://www.beatport.com"

	ctx := context.Background()
	genre := "TRANCE (MAIN FLOOR)"
	songs, err := scraper.FetchSongs(ctx, genre, time.Now())
	if err != nil {
		t.Fatalf("FetchSongs failed: %v", err)
	}

	expectedSongCount := 100
	if len(songs) != expectedSongCount {
		t.Fatalf("Expected to find %d songs, but got %d", expectedSongCount, len(songs))
	}

	// Check the first song
	expectedFirstSong := models.Song{
		Title:  "Children",
		Artist: "ROBERT MILES, TINLICKER",
		Genre:  "TRANCE (MAIN FLOOR)",
		URL:    "https://www.beatport.com/track/track/14269418",
	}

	if songs[0].Title != expectedFirstSong.Title {
		t.Errorf("Expected first song title to be %q, got %q", expectedFirstSong.Title, songs[0].Title)
	}
	if songs[0].Artist != expectedFirstSong.Artist {
		t.Errorf("Expected first song artist to be %q, got %q", expectedFirstSong.Artist, songs[0].Artist)
	}
	if songs[0].Genre != expectedFirstSong.Genre {
		t.Errorf("Expected first song genre to be %q, got %q", expectedFirstSong.Genre, songs[0].Genre)
	}
	if songs[0].URL != expectedFirstSong.URL {
		t.Errorf("Expected first song URL to be %q, got %q", expectedFirstSong.URL, songs[0].URL)
	}

	// Check the second song
	expectedSecondSong := models.Song{
		Title:  "1998 (Victor Ruiz Extended Remix)",
		Artist: "BINARY FINARY",
		Genre:  "TRANCE (MAIN FLOOR)",
		URL:    "https://www.beatport.com/track/track/18160649",
	}

	if songs[1].Title != expectedSecondSong.Title {
		t.Errorf("Expected second song title to be %q, got %q", expectedSecondSong.Title, songs[1].Title)
	}
	if songs[1].Artist != expectedSecondSong.Artist {
		t.Errorf("Expected second song artist to be %q, got %q", expectedSecondSong.Artist, songs[1].Artist)
	}
	if songs[1].Genre != expectedSecondSong.Genre {
		t.Errorf("Expected second song genre to be %q, got %q", expectedSecondSong.Genre, songs[1].Genre)
	}
	if songs[1].URL != expectedSecondSong.URL {
		t.Errorf("Expected second song URL to be %q, got %q", expectedSecondSong.URL, songs[1].URL)
	}
}

func TestBeatstatsScraperFetchSongDetails(t *testing.T) {
	// TODO create a mock server that returns song details HTML
	// TODO Then test the FetchSongDetails method
}
