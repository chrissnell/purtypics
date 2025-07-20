package editor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/deploy"
)

// handleDeploy starts deployment
func (s *Server) handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load deployment config
	deployDir := filepath.Dir(s.MetadataPath)
	config, err := deploy.LoadConfig(deployDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load deployment config: %v", err), http.StatusInternalServerError)
		return
	}

	// Determine deployment type
	var deployType string
	if config.Rsync != nil {
		deployType = "rsync"
	} else if config.S3 != nil {
		deployType = "s3"
	} else if config.Cloudflare != nil {
		deployType = "cloudflare"
	} else {
		http.Error(w, "No deployment configuration found", http.StatusBadRequest)
		return
	}

	// Check if output directory exists
	outputPath := s.OutputPath
	if outputPath == "" {
		outputPath = filepath.Join(s.SourcePath, "output")
	}
	if _, err := os.Stat(outputPath); err != nil {
		http.Error(w, "Output directory not found. Please generate the gallery first.", http.StatusBadRequest)
		return
	}

	// Reset progress
	s.deployTracker.Reset()
	s.deployTracker.Update(0, "starting")

	// Start deployment in background
	go s.runDeployment(config, deployType)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// handleDeployProgress returns deployment progress
func (s *Server) handleDeployProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	progress, status, error := s.deployTracker.Get()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"progress": progress,
		"status":   status,
		"error":    error,
	})
}

// runDeployment runs the deployment process
func (s *Server) runDeployment(config *deploy.Config, deployType string) {
	// Update status
	s.deployTracker.Update(0, "deploying")

	// Determine source path (generated gallery)
	sourcePath := s.OutputPath
	if sourcePath == "" {
		sourcePath = filepath.Join(s.SourcePath, "output")
	}

	// Create deployer based on type
	switch deployType {
	case "rsync":
		if config.Rsync == nil {
			s.deployTracker.SetError("Rsync configuration not found")
			return
		}
		// Create rsync deployer
		deployer := deploy.NewRsyncDeployer(config.Rsync, sourcePath)
		
		// Set progress callback
		deployer.SetProgressCallback(func(progress int, message string) {
			s.deployTracker.Update(progress, message)
		})
		
		// Run deployment
		if err := deployer.Deploy(); err != nil {
			s.deployTracker.SetError(fmt.Sprintf("Deployment failed: %v", err))
			return
		}
		
	case "s3":
		s.deployTracker.SetError("S3 deployment not yet implemented")
		return
		
	case "cloudflare":
		s.deployTracker.SetError("Cloudflare deployment not yet implemented")
		return
		
	default:
		s.deployTracker.SetError(fmt.Sprintf("Unknown deployment type: %s", deployType))
		return
	}

	// Update completion status
	s.deployTracker.Update(100, "completed")
}