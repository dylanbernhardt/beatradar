package models

import (
	"time"
)

// Song represents a single track with its associated metadata
type Song struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	RemixArtist string    `json:"remix_artist,omitempty"`
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"release_date"`
	Label       string    `json:"label"`
	BPM         int       `json:"bpm"`
	Key         string    `json:"key"`
	Length      int       `json:"length"` // in seconds
	URL         string    `json:"url"`    // link to the song on Beatstats
	Artwork     string    `json:"artwork,omitempty"`
	Preview     string    `json:"preview,omitempty"` // URL to audio preview if available
}

// NewSong creates a new Song instance
func NewSong(id, title, artist, genre string) *Song {
	return &Song{
		ID:     id,
		Title:  title,
		Artist: artist,
		Genre:  genre,
	}
}

// SetRemixArtist sets the remix artist for the song
func (s *Song) SetRemixArtist(artist string) {
	s.RemixArtist = artist
}

// SetReleaseInfo sets the release date and label for the song
func (s *Song) SetReleaseInfo(date time.Time, label string) {
	s.ReleaseDate = date
	s.Label = label
}

// SetTechnicalInfo sets the BPM, key, and length for the song
func (s *Song) SetTechnicalInfo(bpm int, key string, length int) {
	s.BPM = bpm
	s.Key = key
	s.Length = length
}

// SetLinks sets the URL and preview link for the song
func (s *Song) SetLinks(url, preview string) {
	s.URL = url
	s.Preview = preview
}

// SetArtwork sets the artwork URL for the song
func (s *Song) SetArtwork(artwork string) {
	s.Artwork = artwork
}
