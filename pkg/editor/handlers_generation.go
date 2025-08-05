package editor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cjs/purtypics/pkg/gallery"
)

// handleGenerate starts gallery generation
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reset progress
	s.genTracker.Reset()
	s.genTracker.Update(0, "running")

	// Start generation in background
	go s.runGeneration()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// handleGenerateProgress returns generation progress
func (s *Server) handleGenerateProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	progress, status, error := s.genTracker.Get()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"progress": progress,
		"status":   status,
		"error":    error,
	})
}

// runGeneration runs the gallery generation process
func (s *Server) runGeneration() {
	// Update status
	s.genTracker.Update(0, "generating")

	// Use the output path from server configuration
	outputPath := s.OutputPath

	// Create generator with progress callback
	generator := gallery.NewGenerator(s.SourcePath, outputPath, s.metadata.Title, "", false)
	generator.MetadataPath = s.MetadataPath
	
	// First, scan to count total photos
	s.genTracker.Update(1, "Scanning albums...")
	
	albums, err := gallery.ScanDirectory(s.SourcePath)
	if err != nil {
		s.genTracker.SetError(fmt.Sprintf("Failed to scan directory: %v", err))
		return
	}
	
	// Count total photos across all albums
	totalPhotos := 0
	for _, album := range albums {
		totalPhotos += len(album.Photos)
	}
	
	s.genTracker.Update(3, fmt.Sprintf("Found %d albums with %d photos", len(albums), totalPhotos))
	
	// Track progress
	var processedPhotos int
	var currentAlbumName string
	
	generator.ProgressCallback = func(current, total int, message string) {
		if strings.HasPrefix(message, "Processing album:") {
			currentAlbumName = strings.TrimPrefix(message, "Processing album: ")
			// Don't update progress here, wait for actual photo processing
		} else if strings.HasPrefix(message, "Processing photos in ") {
			// Starting a new album's photos
			s.genTracker.Update(5 + (85 * processedPhotos / totalPhotos), 
				fmt.Sprintf("Processing %s (%d/%d total photos)", currentAlbumName, processedPhotos, totalPhotos))
		} else if strings.Contains(message, ": ") && current > 0 && total > 0 {
			// Processing individual photos
			processedPhotos++
			progress := 5 + (85 * processedPhotos / totalPhotos)
			if progress > 90 {
				progress = 90
			}
			s.genTracker.Update(progress, fmt.Sprintf("Processing photos (%d/%d)", processedPhotos, totalPhotos))
		} else if message == "Generating HTML pages" {
			s.genTracker.Update(92, "Generating HTML pages...")
		} else if message == "Gallery generation complete" {
			s.genTracker.Update(100, "completed")
		}
	}

	// Run generation
	fmt.Printf("Generating gallery from %s...\n", s.SourcePath)
	if err := generator.Generate(); err != nil {
		s.genTracker.SetError(fmt.Sprintf("Generation failed: %v", err))
		fmt.Printf("Gallery generation failed: %v\n", err)
		return
	}

	// Ensure we reach 100% completion
	s.genTracker.Update(100, "completed")
	fmt.Printf("\nGallery generated successfully at %s\n", outputPath)
}