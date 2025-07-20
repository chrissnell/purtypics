package deploy

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// RsyncDeployer handles rsync deployments
type RsyncDeployer struct {
	config *RsyncConfig
	output string
}

// NewRsyncDeployer creates a new rsync deployer
func NewRsyncDeployer(config *RsyncConfig, outputPath string) *RsyncDeployer {
	return &RsyncDeployer{
		config: config,
		output: outputPath,
	}
}

// Deploy performs the rsync deployment
func (r *RsyncDeployer) Deploy() error {
	if err := r.validate(); err != nil {
		return err
	}

	args := r.buildArgs()
	
	cmd := exec.Command("rsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Executing:", strings.Join(cmd.Args, " "))
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsync failed: %w", err)
	}

	return nil
}

// DeployWithProgress performs the rsync deployment with progress tracking
func (r *RsyncDeployer) DeployWithProgress(progressFn func(int)) error {
	if err := r.validate(); err != nil {
		return err
	}

	args := r.buildArgs()
	
	cmd := exec.Command("rsync", args...)
	
	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe: %w", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting rsync: %w", err)
	}
	
	// Pattern to match rsync progress output
	// Example: "1,234,567  32%   12.34MB/s    0:00:12"
	progressRegex := regexp.MustCompile(`^\s*[\d,]+\s+(\d+)%`)
	
	// Read stdout and parse progress
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Still print to console
			
			// Try to extract progress percentage
			if matches := progressRegex.FindStringSubmatch(line); len(matches) > 1 {
				if percent, err := strconv.Atoi(matches[1]); err == nil {
					progressFn(percent)
				}
			}
		}
	}()
	
	// Read stderr and print errors
	go func() {
		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			fmt.Fprintln(os.Stderr, errScanner.Text())
		}
	}()
	
	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("rsync failed: %w", err)
	}
	
	// Ensure we report 100% completion
	progressFn(100)
	
	return nil
}

// validate checks if the configuration is valid
func (r *RsyncDeployer) validate() error {
	if r.config.Host == "" {
		return fmt.Errorf("rsync host is required")
	}
	if r.config.Path == "" {
		return fmt.Errorf("rsync path is required")
	}
	
	// Check if rsync is available
	if _, err := exec.LookPath("rsync"); err != nil {
		return fmt.Errorf("rsync not found in PATH")
	}

	// Check if output directory exists
	if _, err := os.Stat(r.output); err != nil {
		return fmt.Errorf("output directory not found: %s", r.output)
	}

	return nil
}

// buildArgs builds the rsync command arguments
func (r *RsyncDeployer) buildArgs() []string {
	args := []string{
		"-av",               // archive, verbose (no compression for images)
		"--progress",        // show progress
		"--human-readable",  // human readable sizes
		"--update",          // skip files that are newer on receiver
		"--checksum",        // skip based on checksum, not mod-time & size
	}

	// Add port if specified
	if r.config.Port != 0 && r.config.Port != 22 {
		args = append(args, "-e", fmt.Sprintf("ssh -p %d", r.config.Port))
	}

	// Add dry run flag if specified
	if r.config.DryRun {
		args = append(args, "--dry-run")
	}

	// Add source (trailing slash to copy contents)
	args = append(args, r.output+"/")

	// Add destination
	dest := fmt.Sprintf("%s:%s", r.config.Host, r.config.Path)
	args = append(args, dest)

	return args
}

// GetInfo returns deployment information
func (r *RsyncDeployer) GetInfo() string {
	return fmt.Sprintf("rsync to %s:%s", r.config.Host, r.config.Path)
}