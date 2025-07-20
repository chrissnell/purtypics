package editor

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cjs/purtypics/pkg/gallery"
)

// handleAlbums returns album information
func (s *Server) handleAlbums(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type albumResponse struct {
		Path        string    `json:"path"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		PhotoCount  int       `json:"photoCount"`
		CoverPhoto  string    `json:"coverPhoto"`
		Hidden      bool      `json:"hidden"`
		Date        time.Time `json:"date"`
		Photos      []string  `json:"photos"`
	}

	albums := make([]albumResponse, len(s.albums))
	for i, album := range s.albums {
		// Get photo filenames
		photos := make([]string, len(album.Photos))
		for j, photo := range album.Photos {
			photos[j] = photo.Filename
		}
		
		resp := albumResponse{
			Path:       album.Path,
			Title:      album.Title,
			PhotoCount: len(album.Photos),
			Photos:     photos,
		}

		// Apply metadata if exists
		if meta := s.metadata.GetAlbumMetadata(album.Path); meta != nil {
			resp.Title = meta.Title
			resp.Description = meta.Description
			resp.CoverPhoto = meta.CoverPhoto
			resp.Hidden = meta.Hidden
			resp.Date = meta.Date
		}

		albums[i] = resp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

// handlePhotos returns photos for a specific album
func (s *Server) handlePhotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract album path from URL
	albumPath := strings.TrimPrefix(r.URL.Path, "/api/photos/")
	albumPath = filepath.Join(s.SourcePath, albumPath)

	// Find the album
	var album *gallery.Album
	for i := range s.albums {
		if s.albums[i].Path == albumPath {
			album = &s.albums[i]
			break
		}
	}

	if album == nil {
		http.Error(w, "Album not found", http.StatusNotFound)
		return
	}

	type photoResponse struct {
		Path        string            `json:"path"`
		Filename    string            `json:"filename"`
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Hidden      bool              `json:"hidden"`
		IsVideo     bool              `json:"isVideo"`
		Thumbnails  map[string]string `json:"thumbnails"`
	}

	photos := make([]photoResponse, len(album.Photos))
	for i, photo := range album.Photos {
		photoPath := filepath.Join(album.Path, photo.Filename)
		resp := photoResponse{
			Path:       photoPath,
			Filename:   photo.Filename,
			Title:      photo.Title,
			IsVideo:    photo.IsVideo,
			Thumbnails: photo.Thumbnails,
		}

		// Apply metadata if exists
		if meta := s.metadata.GetPhotoMetadata(photoPath); meta != nil {
			resp.Title = meta.Title
			resp.Description = meta.Description
			resp.Hidden = meta.Hidden
		}

		photos[i] = resp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}