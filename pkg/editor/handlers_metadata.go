package editor

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/deploy"
	"github.com/cjs/purtypics/pkg/metadata"
)

// handleMetadata returns the current metadata
func (s *Server) handleMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.metadata)
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

// handleDeployConfig returns deployment configuration
func (s *Server) handleDeployConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Check if deploy.yaml exists
		deployDir := filepath.Dir(s.MetadataPath)
		config, err := deploy.LoadConfig(deployDir)
		if err != nil {
			// Return empty config if file doesn't exist
			config = &deploy.Config{}
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
		
	case http.MethodPost:
		var config deploy.Config
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Save config
		deployDir := filepath.Dir(s.MetadataPath)
		if err := deploy.SaveConfig(deployDir, &config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}