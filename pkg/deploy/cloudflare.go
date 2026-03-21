package deploy

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const cloudflareAPIBase = "https://api.cloudflare.com/client/v4"

// CloudflareDeployer handles deployments to Cloudflare Pages.
type CloudflareDeployer struct {
	config           *CloudflareConfig
	output           string
	token            string
	progressCallback func(int, string)
}

// cloudflareResponse is the envelope for all Cloudflare API responses.
type cloudflareResponse struct {
	Success  bool                 `json:"success"`
	Errors   []cloudflareError    `json:"errors"`
	Messages []cloudflareMessage  `json:"messages"`
	Result   json.RawMessage      `json:"result"`
}

type cloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cloudflareMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type deploymentResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// NewCloudflareDeployer creates a new Cloudflare Pages deployer.
// The API token is read from the CLOUDFLARE_API_TOKEN environment variable.
func NewCloudflareDeployer(config *CloudflareConfig, outputPath string) (*CloudflareDeployer, error) {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_TOKEN environment variable is not set")
	}
	return &CloudflareDeployer{
		config: config,
		output: outputPath,
		token:  token,
	}, nil
}

// SetProgressCallback sets the progress callback function.
func (c *CloudflareDeployer) SetProgressCallback(fn func(int, string)) {
	c.progressCallback = fn
}

// TestConnection verifies that the API token is valid and the project exists.
func (c *CloudflareDeployer) TestConnection() error {
	if err := c.validate(); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/accounts/%s/pages/projects/%s",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to Cloudflare API: %w", err)
	}
	defer resp.Body.Close()

	var cfResp cloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return fmt.Errorf("decoding Cloudflare response: %w", err)
	}

	if !cfResp.Success {
		return fmt.Errorf("Cloudflare API error: %s", formatErrors(cfResp.Errors))
	}

	return nil
}

// Deploy uploads the output directory to Cloudflare Pages as a new deployment.
func (c *CloudflareDeployer) Deploy() error {
	if err := c.validate(); err != nil {
		return err
	}

	c.progress(0, "Collecting files...")

	// Collect all files to upload.
	files, err := c.collectFiles()
	if err != nil {
		return fmt.Errorf("collecting files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found in output directory: %s", c.output)
	}

	c.progress(5, fmt.Sprintf("Uploading %d files...", len(files)))

	// Build multipart body using a pipe so we don't buffer everything in memory.
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Write multipart parts in a goroutine.
	errCh := make(chan error, 1)
	go func() {
		defer pw.Close()
		defer writer.Close()

		for i, relPath := range files {
			absPath := filepath.Join(c.output, relPath)
			f, err := os.Open(absPath)
			if err != nil {
				errCh <- fmt.Errorf("opening %s: %w", relPath, err)
				return
			}

			// Cloudflare expects the part name to be the path with leading slash.
			partName := "/" + relPath
			part, err := writer.CreateFormFile(partName, filepath.Base(relPath))
			if err != nil {
				f.Close()
				errCh <- fmt.Errorf("creating form part for %s: %w", relPath, err)
				return
			}

			if _, err := io.Copy(part, f); err != nil {
				f.Close()
				errCh <- fmt.Errorf("writing %s: %w", relPath, err)
				return
			}
			f.Close()

			// Report progress (5% to 90% range for uploads).
			pct := 5 + (85 * (i + 1) / len(files))
			c.progress(pct, fmt.Sprintf("Uploading %d/%d files...", i+1, len(files)))
		}

		// Add branch if configured.
		if c.config.Branch != "" {
			if err := writer.WriteField("branch", c.config.Branch); err != nil {
				errCh <- fmt.Errorf("writing branch field: %w", err)
				return
			}
		}

		errCh <- nil
	}()

	// Send the request.
	url := fmt.Sprintf("%s/accounts/%s/pages/projects/%s/deployments",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	req, err := http.NewRequest(http.MethodPost, url, pr)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	c.progress(90, "Waiting for Cloudflare...")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Drain the writer goroutine error too.
		<-errCh
		return fmt.Errorf("uploading to Cloudflare: %w", err)
	}
	defer resp.Body.Close()

	// Check the writer goroutine for errors.
	if writeErr := <-errCh; writeErr != nil {
		return writeErr
	}

	var cfResp cloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return fmt.Errorf("decoding Cloudflare response: %w", err)
	}

	if !cfResp.Success {
		return fmt.Errorf("deployment failed: %s", formatErrors(cfResp.Errors))
	}

	// Extract deployment URL.
	var result deploymentResult
	if err := json.Unmarshal(cfResp.Result, &result); err == nil && result.URL != "" {
		c.progress(100, fmt.Sprintf("Deployed to %s", result.URL))
		fmt.Printf("Deployment URL: %s\n", result.URL)
	} else {
		c.progress(100, "Deployment complete")
	}

	return nil
}

// GetInfo returns a human-readable description of the deployment target.
func (c *CloudflareDeployer) GetInfo() string {
	return fmt.Sprintf("Cloudflare Pages project %q", c.config.Project)
}

func (c *CloudflareDeployer) validate() error {
	if c.config.Project == "" {
		return fmt.Errorf("Cloudflare Pages project name is required")
	}
	if c.config.AccountID == "" {
		return fmt.Errorf("Cloudflare account ID is required")
	}
	if _, err := os.Stat(c.output); err != nil {
		return fmt.Errorf("output directory not found: %s", c.output)
	}
	return nil
}

// collectFiles walks the output directory and returns relative paths for all files,
// excluding config files that shouldn't be deployed.
func (c *CloudflareDeployer) collectFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(c.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(c.output, path)
		if err != nil {
			return err
		}

		// Use forward slashes for the API.
		rel = filepath.ToSlash(rel)

		// Skip config files.
		if rel == "gallery.yaml" || rel == "deploy.yaml" {
			return nil
		}

		files = append(files, rel)
		return nil
	})
	return files, err
}

func (c *CloudflareDeployer) progress(pct int, msg string) {
	if c.progressCallback != nil {
		c.progressCallback(pct, msg)
	}
}

func formatErrors(errs []cloudflareError) string {
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Message
	}
	return strings.Join(msgs, "; ")
}
