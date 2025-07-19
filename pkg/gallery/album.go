package gallery

import (
	"os"
	"path/filepath"
	"sort"
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

// GetNewestPhotoTime returns the most recent photo time in the album
func (a *Album) GetNewestPhotoTime() time.Time {
	var newest time.Time
	for _, photo := range a.Photos {
		if photo.EXIF != nil && !photo.EXIF.DateTime.IsZero() {
			if newest.IsZero() || photo.EXIF.DateTime.After(newest) {
				newest = photo.EXIF.DateTime
			}
		}
	}
	// If no EXIF dates found, use current time
	if newest.IsZero() {
		newest = time.Now()
	}
	return newest
}

// SortPhotosByDate sorts photos by EXIF date (oldest to newest)
func (a *Album) SortPhotosByDate() {
	sort.Slice(a.Photos, func(i, j int) bool {
		timeI := time.Now()
		timeJ := time.Now()
		
		if a.Photos[i].EXIF != nil && !a.Photos[i].EXIF.DateTime.IsZero() {
			timeI = a.Photos[i].EXIF.DateTime
		}
		if a.Photos[j].EXIF != nil && !a.Photos[j].EXIF.DateTime.IsZero() {
			timeJ = a.Photos[j].EXIF.DateTime
		}
		
		return timeI.Before(timeJ)
	})
}

// SortAlbumsByDate sorts albums by their newest photo (newest album first)
func SortAlbumsByDate(albums []Album) {
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].GetNewestPhotoTime().After(albums[j].GetNewestPhotoTime())
	})
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