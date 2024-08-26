package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBeatstatsScraperFetchSongs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO mock HTML that your scraper expects
		w.Write([]byte(`
            <div class="track-row">
                <div class="track-title"><a href="/song/1">Test Song</a></div>
                <div class="track-artists">Test Artist</div>
            </div>
        `))
	}))
	defer server.Close()

	scraper := NewBeatstatsScraper(server.URL)

	songs, err := scraper.FetchSongs(context.Background(), "House", time.Now())
	if err != nil {
		t.Fatalf("FetchSongs failed: %v", err)
	}

	if len(songs) != 1 {
		t.Fatalf("Expected 1 song, got %d", len(songs))
	}

	if songs[0].Title != "Test Song" {
		t.Errorf("Expected song title 'Test Song', got '%s'", songs[0].Title)
	}

	if songs[0].Artist != "Test Artist" {
		t.Errorf("Expected artist 'Test Artist', got '%s'", songs[0].Artist)
	}
}

func TestBeatstatsScraperFetchSongDetails(t *testing.T) {
	// TODO create a mock server that returns song details HTML
	// TODO Then test the FetchSongDetails method
}
