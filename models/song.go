package models

import (
	"time"
)

// Song holds information about song that is used to
// make recommendations.
type Song struct {
	ID            string
	Title         string
	Author        string
	SongURL       string
	ImageURL      string
	Description   string
	Recommended   bool
	CreatedAt     time.Time
	RecommendedAt time.Time
}
