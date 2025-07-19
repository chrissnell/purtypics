package editor

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cjs/purtypics/pkg/gallery"
	"github.com/cjs/purtypics/pkg/metadata"
	"github.com/disintegration/imaging"
)

// Server provides a web interface for editing gallery metadata
type Server struct {
	SourcePath   string
	MetadataPath string
	Port         int
	metadata     *metadata.GalleryMetadata
	albums       []gallery.Album
}

// NewServer creates a new editor server
func NewServer(sourcePath, metadataPath string, port int) *Server {
	if metadataPath == "" {
		metadataPath = filepath.Join(sourcePath, "gallery.yaml")
	}
	
	return &Server{
		SourcePath:   sourcePath,
		MetadataPath: metadataPath,
		Port:         port,
	}
}

// Start runs the web server
func (s *Server) Start() error {
	// Load existing metadata
	meta, err := metadata.Load(s.MetadataPath)
	if err != nil {
		return fmt.Errorf("loading metadata: %w", err)
	}
	s.metadata = meta

	// Scan albums
	albums, err := gallery.ScanDirectory(s.SourcePath)
	if err != nil {
		return fmt.Errorf("scanning albums: %w", err)
	}
	s.albums = albums

	// Set up routes
	mux := http.NewServeMux()
	
	// API routes
	mux.HandleFunc("/api/metadata", s.handleMetadata)
	mux.HandleFunc("/api/albums", s.handleAlbums)
	mux.HandleFunc("/api/photos/", s.handlePhotos)
	mux.HandleFunc("/api/save", s.handleSave)
	
	// Static files
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static/", s.handleStatic)
	
	// Serve images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(s.SourcePath))))
	
	// Serve thumbnails (dynamic generation)
	mux.HandleFunc("/thumbs/", s.handleThumbnails)

	// Find available port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		// If the specified port is in use, find an available one
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return fmt.Errorf("failed to find available port: %w", err)
		}
	}
	
	// Get the actual port
	addr := listener.Addr().(*net.TCPAddr)
	s.Port = addr.Port

	fmt.Printf("Starting metadata editor on http://localhost:%d\n", s.Port)
	
	// Serve on the listener
	return http.Serve(listener, mux)
}

// handleMetadata returns the current metadata
func (s *Server) handleMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.metadata)
}

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
	}

	albums := make([]albumResponse, len(s.albums))
	for i, album := range s.albums {
		resp := albumResponse{
			Path:       album.Path,
			Title:      album.Title,
			PhotoCount: len(album.Photos),
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

// handleSave saves the updated metadata
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var meta metadata.GalleryMetadata
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save to file
	if err := metadata.Save(&meta, s.MetadataPath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update in-memory copy
	s.metadata = &meta

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

// handleIndex serves the main editor HTML
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(editorHTML))
}

// handleStatic serves static assets (CSS, JS)
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	
	switch path {
	case "editor.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(editorCSS))
	case "editor.js":
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(editorJS))
	default:
		http.NotFound(w, r)
	}
}

// handleThumbnails generates thumbnails on the fly
func (s *Server) handleThumbnails(w http.ResponseWriter, r *http.Request) {
	// Parse the URL: /thumbs/{size}/{album}/{photo}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/thumbs/"), "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	
	size := parts[0]
	album := parts[1]
	photo := parts[2]
	
	// Validate size
	var width, height int
	switch size {
	case "small":
		width, height = 300, 200
	case "medium":
		width, height = 600, 400
	case "large":
		width, height = 1200, 800
	default:
		http.NotFound(w, r)
		return
	}
	
	// Build source path
	sourcePath := filepath.Join(s.SourcePath, album, photo)
	
	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	
	// Open source image
	file, err := os.Open(sourcePath)
	if err != nil {
		http.Error(w, "Error opening image", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	
	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Error decoding image", http.StatusInternalServerError)
		return
	}
	
	// Create thumbnail - fit within bounds while maintaining aspect ratio
	thumb := imaging.Fit(img, width, height, imaging.Lanczos)
	
	// Set content type and encode
	switch format {
	case "jpeg", "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
		jpeg.Encode(w, thumb, &jpeg.Options{Quality: 85})
	case "png":
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		png.Encode(w, thumb)
	default:
		// Default to JPEG
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		jpeg.Encode(w, thumb, &jpeg.Options{Quality: 85})
	}
}