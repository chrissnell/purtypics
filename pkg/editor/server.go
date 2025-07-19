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
	"sync"
	"time"

	"github.com/cjs/purtypics/pkg/exif"
	"github.com/cjs/purtypics/pkg/gallery"
	"github.com/cjs/purtypics/pkg/metadata"
	"github.com/disintegration/imaging"
)

// Server provides a web interface for editing gallery metadata
type Server struct {
	SourcePath   string
	MetadataPath string
	OutputPath   string
	Port         int
	metadata     *metadata.GalleryMetadata
	albums       []gallery.Album
	
	// Generation progress tracking
	genProgress   int
	genStatus     string
	genError      string
	genMutex      sync.RWMutex
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

// GetActualPort finds an available port and returns it
func (s *Server) GetActualPort() (int, net.Listener, error) {
	// Find available port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		// If the specified port is in use, find an available one
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return 0, nil, fmt.Errorf("failed to find available port: %w", err)
		}
	}
	
	// Get the actual port
	addr := listener.Addr().(*net.TCPAddr)
	s.Port = addr.Port
	return s.Port, listener, nil
}

// Start runs the web server
func (s *Server) Start() error {
	return s.StartWithListener(nil)
}

// StartWithListener runs the web server with a provided listener
func (s *Server) StartWithListener(listener net.Listener) error {
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
	mux.HandleFunc("/api/generate", s.handleGenerate)
	mux.HandleFunc("/api/generate/progress", s.handleGenerateProgress)
	
	// Static files
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static/", s.handleStatic)
	
	// Serve images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(s.SourcePath))))
	
	// Serve thumbnails (dynamic generation)
	mux.HandleFunc("/thumbs/", s.handleThumbnails)

	// If no listener provided, create one
	if listener == nil {
		var err error
		_, listener, err = s.GetActualPort()
		if err != nil {
			return err
		}
	}
	
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
	
	// Apply EXIF orientation manually
	orientation, err := exif.GetOrientation(sourcePath)
	if err != nil {
		fmt.Printf("Error getting orientation for %s: %v\n", sourcePath, err)
	} else {
		fmt.Printf("Applying orientation %d to %s\n", orientation, photo)
	}
	img = applyOrientation(img, orientation)
	
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

// handleGenerate starts the gallery generation process
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check if output path is set
	if s.OutputPath == "" {
		http.Error(w, "Output path not configured", http.StatusInternalServerError)
		return
	}
	
	// Reset progress
	s.genMutex.Lock()
	s.genProgress = 0
	s.genStatus = "running"
	s.genError = ""
	s.genMutex.Unlock()
	
	// Start generation in background
	go s.runGeneration()
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// handleGenerateProgress returns the current generation progress
func (s *Server) handleGenerateProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	s.genMutex.RLock()
	progress := s.genProgress
	status := s.genStatus
	genError := s.genError
	s.genMutex.RUnlock()
	
	response := map[string]interface{}{
		"progress": progress,
		"status":   status,
	}
	
	if genError != "" {
		response["error"] = genError
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// runGeneration performs the actual gallery generation
func (s *Server) runGeneration() {
	// Create a progress callback
	progressCallback := func(current, total int, message string) {
		s.genMutex.Lock()
		if total > 0 {
			s.genProgress = (current * 100) / total
		}
		s.genMutex.Unlock()
	}
	
	// Create generator with progress callback
	generator := gallery.NewGenerator(s.SourcePath, s.OutputPath, s.metadata.Title, "", false)
	generator.MetadataPath = s.MetadataPath
	generator.ProgressCallback = progressCallback
	
	// Run generation
	if err := generator.Generate(); err != nil {
		s.genMutex.Lock()
		s.genStatus = "error"
		s.genError = err.Error()
		s.genMutex.Unlock()
		return
	}
	
	// Mark as completed
	s.genMutex.Lock()
	s.genProgress = 100
	s.genStatus = "completed"
	s.genMutex.Unlock()
}

// applyOrientation applies EXIF orientation transformations to an image
func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		// Normal - no rotation needed
		return img
	case 2:
		// Flip horizontal
		return imaging.FlipH(img)
	case 3:
		// Rotate 180 degrees
		return imaging.Rotate180(img)
	case 4:
		// Flip vertical
		return imaging.FlipV(img)
	case 5:
		// Rotate 90 degrees clockwise and flip horizontally
		return imaging.FlipH(imaging.Rotate90(img))
	case 6:
		// Rotate 90 degrees clockwise (actually 270 in imaging library)
		return imaging.Rotate270(img)
	case 7:
		// Rotate 90 degrees counterclockwise and flip horizontally
		return imaging.FlipH(imaging.Rotate270(img))
	case 8:
		// Rotate 90 degrees counterclockwise (actually 90 in imaging library)
		return imaging.Rotate90(img)
	default:
		// Unknown orientation, return as-is
		return img
	}
}