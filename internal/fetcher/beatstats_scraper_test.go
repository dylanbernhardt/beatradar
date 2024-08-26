package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewBeatstatsScraper(server.URL)

	ctx := context.Background()
	songs, err := scraper.FetchSongs(ctx, "TRANCE (MAIN FLOOR)", time.Now())
	if err != nil {
		t.Fatalf("FetchSongs failed: %v", err)
	}

	if len(songs) == 0 {
		t.Fatalf("Expected to find songs, but got none")
	}

	expectedFirstSong := models.Song{
		Title:  "Children",
		Artist: "ROBERT MILES, TINLICKER",
		Genre:  "TRANCE (MAIN FLOOR)",
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

	expectedSecondSong := models.Song{
		Title:  "1998 (Victor Ruiz Extended Remix)",
		Artist: "BINARY FINARY",
		Genre:  "TRANCE (MAIN FLOOR)",
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

	for i, song := range songs {
		if !strings.HasPrefix(song.URL, "https://www.beatport.com/track/") {
			t.Errorf("Song %d doesn't have a valid Beatport URL: %s", i+1, song.URL)
		}
	}
}

func TestBeatstatsScraperFetchSongDetails(t *testing.T) {
	// TODO create a mock server that returns song details HTML
	// TODO Then test the FetchSongDetails method
}
