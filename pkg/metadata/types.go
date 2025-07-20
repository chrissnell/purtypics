package metadata

import (
	"time"
)

// GalleryMetadata represents the overall gallery configuration
type GalleryMetadata struct {
	Title         string                     `yaml:"title" json:"title"`
	Description   string                     `yaml:"description" json:"description"`
	Author        string                     `yaml:"author" json:"author"`
	Copyright     string                     `yaml:"copyright" json:"copyright"`
	ShowLocations bool                       `yaml:"show_locations" json:"show_locations"`
	Albums        map[string]*AlbumMetadata  `yaml:"albums" json:"albums"`
	Photos        map[string]*PhotoMetadata  `yaml:"photos" json:"photos"`
}

// AlbumMetadata represents metadata for a single album
type AlbumMetadata struct {
	Title       string    `yaml:"title" json:"title"`
	Description string    `yaml:"description" json:"description"`
	Date        time.Time `yaml:"date" json:"date"`
	CoverPhoto  string    `yaml:"cover_photo" json:"cover_photo"`
	Hidden      bool      `yaml:"hidden" json:"hidden"`
	SortOrder   string    `yaml:"sort_order" json:"sort_order"` // "date", "name", "custom"
	CustomOrder []string  `yaml:"custom_order" json:"custom_order"`
	Tags        []string  `yaml:"tags" json:"tags"`
}

// PhotoMetadata represents metadata for a single photo
type PhotoMetadata struct {
	Title       string   `yaml:"title" json:"title"`
	Description string   `yaml:"description" json:"description"`
	Tags        []string `yaml:"tags" json:"tags"`
	Hidden      bool     `yaml:"hidden" json:"hidden"`
	SortIndex   int      `yaml:"sort_index" json:"sort_index"`
}