package models

import (
	"testing"
	"time"
)

func TestNewSong(t *testing.T) {
	song := NewSong("1", "Test Song", "Test Artist", "House")

	if song.ID != "1" {
		t.Errorf("Expected song ID to be '1', got '%s'", song.ID)
	}
	if song.Title != "Test Song" {
		t.Errorf("Expected song title to be 'Test Song', got '%s'", song.Title)
	}
	if song.Artist != "Test Artist" {
		t.Errorf("Expected song artist to be 'Test Artist', got '%s'", song.Artist)
	}
	if song.Genre != "House" {
		t.Errorf("Expected song genre to be 'House', got '%s'", song.Genre)
	}
}

func TestSongMethods(t *testing.T) {
	song := NewSong("1", "Test Song", "Test Artist", "House")

	song.SetRemixArtist("Remix Artist")
	if song.RemixArtist != "Remix Artist" {
		t.Errorf("Expected remix artist to be 'Remix Artist', got '%s'", song.RemixArtist)
	}

	releaseDate := time.Now()
	song.SetReleaseInfo(releaseDate, "Test Label")
	if song.ReleaseDate != releaseDate {
		t.Errorf("Expected release date to be %v, got %v", releaseDate, song.ReleaseDate)
	}
	if song.Label != "Test Label" {
		t.Errorf("Expected label to be 'Test Label', got '%s'", song.Label)
	}

	song.SetTechnicalInfo(128, "Cmaj", 180)
	if song.BPM != 128 {
		t.Errorf("Expected BPM to be 128, got %d", song.BPM)
	}
	if song.Key != "Cmaj" {
		t.Errorf("Expected key to be 'Cmaj', got '%s'", song.Key)
	}
	if song.Length != 180 {
		t.Errorf("Expected length to be 180, got %d", song.Length)
	}

	song.SetLinks("https://example.com", "https://preview.com")
	if song.URL != "https://example.com" {
		t.Errorf("Expected URL to be 'https://example.com', got '%s'", song.URL)
	}
	if song.Preview != "https://preview.com" {
		t.Errorf("Expected preview to be 'https://preview.com', got '%s'", song.Preview)
	}

	song.SetArtwork("https://artwork.com")
	if song.Artwork != "https://artwork.com" {
		t.Errorf("Expected artwork to be 'https://artwork.com', got '%s'", song.Artwork)
	}
}
