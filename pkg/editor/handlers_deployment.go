package editor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/deploy"
)

type deployRequest struct {
	Target string `json:"target"`
	DryRun bool   `json:"dry_run"`
}

// handleDeploy starts deployment
func (s *Server) handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req deployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Load deployment config
	deployDir := filepath.Dir(s.MetadataPath)
	config, err := deploy.LoadConfig(deployDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load deployment config: %v", err), http.StatusInternalServerError)
		return
	}

	// Determine deployment type from request or auto-detect.
	deployType := req.Target
	if deployType == "" {
		switch {
		case config.Rsync != nil:
			deployType = "rsync"
		case config.S3 != nil:
			deployType = "s3"
		case config.Cloudflare != nil:
			deployType = "cloudflare"
		default:
			http.Error(w, "No deployment configuration found", http.StatusBadRequest)
			return
		}
	}

	fmt.Printf("Deployment type: %s (dry_run=%v)\n", deployType, req.DryRun)

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
	go s.runDeployment(config, deployType, req.DryRun)

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
func (s *Server) runDeployment(config *deploy.Config, deployType string, dryRun bool) {
	s.deployTracker.Update(0, "deploying")

	sourcePath := s.OutputPath
	if sourcePath == "" {
		sourcePath = filepath.Join(s.SourcePath, "output")
	}

	fmt.Printf("Starting deployment from %s...\n", sourcePath)

	switch deployType {
	case "rsync":
		if config.Rsync == nil {
			s.deployTracker.SetError("Rsync configuration not found")
			return
		}
		config.Rsync.DryRun = dryRun
		deployer := deploy.NewRsyncDeployer(config.Rsync, sourcePath)
		deployer.SetProgressCallback(func(progress int, message string) {
			s.deployTracker.Update(progress, message)
		})
		if err := deployer.Deploy(); err != nil {
			s.deployTracker.SetError(fmt.Sprintf("Deployment failed: %v", err))
			fmt.Printf("Deployment failed: %v\n", err)
			return
		}

	case "cloudflare":
		if config.Cloudflare == nil {
			s.deployTracker.SetError("Cloudflare configuration not found")
			return
		}
		deployer, err := deploy.NewCloudflareDeployer(config.Cloudflare, sourcePath)
		if err != nil {
			s.deployTracker.SetError(err.Error())
			return
		}
		deployer.SetProgressCallback(func(progress int, message string) {
			s.deployTracker.Update(progress, message)
		})
		if dryRun {
			if err := deployer.TestConnection(); err != nil {
				s.deployTracker.SetError(fmt.Sprintf("Connection test failed: %v", err))
				return
			}
		} else {
			if err := deployer.Deploy(); err != nil {
				s.deployTracker.SetError(fmt.Sprintf("Deployment failed: %v", err))
				fmt.Printf("Deployment failed: %v\n", err)
				return
			}
		}

	case "s3":
		if config.S3 == nil {
			s.deployTracker.SetError("S3 configuration not found")
			return
		}
		s3Deployer := deploy.NewS3Deployer(config.S3, sourcePath)
		s3Deployer.SetProgressCallback(func(progress int, message string) {
			s.deployTracker.Update(progress, message)
		})
		if dryRun {
			if err := s3Deployer.TestConnection(); err != nil {
				s.deployTracker.SetError(fmt.Sprintf("Connection test failed: %v", err))
				return
			}
		} else {
			if err := s3Deployer.Deploy(); err != nil {
				s.deployTracker.SetError(fmt.Sprintf("Deployment failed: %v", err))
				fmt.Printf("Deployment failed: %v\n", err)
				return
			}
		}

	default:
		s.deployTracker.SetError(fmt.Sprintf("Unknown deployment type: %s", deployType))
		return
	}

	s.deployTracker.Update(100, "completed")
	fmt.Printf("\nDeployment completed successfully!\n")
}
