package gallery

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cjs/purtypics/pkg/exif"
)

// Album represents a photo album
type Album struct {
	ID          string
	Title       string
	Description string
	Path        string
	Photos      []Photo
	Thumbnail   string // photo ID for album thumbnail
	CoverPhoto  string // custom cover photo from metadata
	CreatedAt   time.Time
}

// Photo represents a single photo
type Photo struct {
	ID          string
	Filename    string
	Title       string
	Description string
	Path        string
	Width       int
	Height      int
	AspectRatio string
	EXIF        *exif.EXIFData
	Thumbnails  map[string]string // size -> path
	IsVideo     bool
	VideoPath   string // Original video path
}

// supportedFormats lists all supported image and video formats
var supportedFormats = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".heif": true,
	".heic": true,
	".mov":  true,
	".mp4":  true,
	".avi":  true,
}

// ScanDirectory finds all albums in a directory
func ScanDirectory(rootPath string) ([]Album, error) {
	var albums []Album

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		albumPath := filepath.Join(rootPath, entry.Name())
		album, err := scanAlbum(albumPath)
		if err != nil {
			continue // skip invalid albums
		}

		if len(album.Photos) > 0 {
			albums = append(albums, album)
		}
	}

	return albums, nil
}

func scanAlbum(path string) (Album, error) {
	album := Album{
		ID:        filepath.Base(path),
		Title:     formatTitle(filepath.Base(path)),
		Path:      path,
		CreatedAt: time.Now(),
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return album, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !supportedFormats[ext] {
			continue
		}

		photo := Photo{
			ID:       strings.TrimSuffix(entry.Name(), ext),
			Filename: entry.Name(),
			Title:    strings.TrimSuffix(entry.Name(), ext),
			Path:     filepath.Join(path, entry.Name()),
			IsVideo:  isVideoFormat(ext),
		}
		
		// Store original video path
		if photo.IsVideo {
			photo.VideoPath = photo.Path
		}

		album.Photos = append(album.Photos, photo)
	}

	// Use first photo as default thumbnail
	if len(album.Photos) > 0 {
		album.Thumbnail = album.Photos[0].ID
	}

	return album, nil
}

// formatTitle converts directory names to readable titles
func formatTitle(name string) string {
	// Replace underscores with spaces
	title := strings.ReplaceAll(name, "_", " ")
	// Basic title case (could be improved)
	return title
}


// isVideoFormat checks if the extension is a video format
func isVideoFormat(ext string) bool {
	videoFormats := map[string]bool{
		".mov": true,
		".mp4": true,
		".avi": true,
		".webm": true,
		".mkv": true,
	}
	return videoFormats[strings.ToLower(ext)]
}