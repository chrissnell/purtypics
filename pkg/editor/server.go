package editor

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/cjs/purtypics/pkg/gallery"
	"github.com/cjs/purtypics/pkg/metadata"
)

// Server provides a web interface for editing gallery metadata
type Server struct {
	SourcePath   string
	MetadataPath string
	OutputPath   string
	Port         int
	metadata     *metadata.GalleryMetadata
	albums       []gallery.Album
	
	// Progress trackers
	genTracker    *ProgressTracker
	deployTracker *ProgressTracker
	
	// Legacy progress tracking (for compatibility)
	genProgress   int
	genStatus     string
	genError      string
	genMutex      sync.RWMutex
	deployProgress int
	deployStatus   string
	deployError    string
	deployMutex    sync.RWMutex
}

// NewServer creates a new editor server
func NewServer(sourcePath, metadataPath string, port int) *Server {
	if metadataPath == "" {
		metadataPath = filepath.Join(sourcePath, "gallery.yaml")
	}
	
	return &Server{
		SourcePath:    sourcePath,
		MetadataPath:  metadataPath,
		Port:          port,
		genTracker:    NewProgressTracker(),
		deployTracker: NewProgressTracker(),
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
	mux.HandleFunc("/api/deploy-config", s.handleDeployConfig)
	mux.HandleFunc("/api/deploy", s.handleDeploy)
	mux.HandleFunc("/api/deploy/progress", s.handleDeployProgress)
	
	// Static files
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static/", s.handleStatic)
	
	// Serve images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(s.SourcePath))))
	
	// Serve generated gallery
	if s.OutputPath != "" {
		mux.Handle("/gallery/", http.StripPrefix("/gallery/", http.FileServer(http.Dir(s.OutputPath))))
	}
	
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
