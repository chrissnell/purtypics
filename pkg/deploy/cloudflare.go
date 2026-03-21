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
// If AutoCreate is enabled and the project doesn't exist, it will be created.
func (c *CloudflareDeployer) TestConnection() error {
	if err := c.validate(); err != nil {
		return err
	}
	return c.ensureProject()
}

// ensureProject checks if the project exists and creates it if AutoCreate is set.
func (c *CloudflareDeployer) ensureProject() error {
	exists, err := c.projectExists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if !c.config.AutoCreate {
		return fmt.Errorf("project %q not found (enable auto-create to create it)", c.config.Project)
	}
	return c.createProject()
}

// projectExists checks whether the Pages project exists.
func (c *CloudflareDeployer) projectExists() (bool, error) {
	url := fmt.Sprintf("%s/accounts/%s/pages/projects/%s",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("connecting to Cloudflare API: %w", err)
	}
	defer resp.Body.Close()

	var cfResp cloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return false, fmt.Errorf("decoding Cloudflare response: %w", err)
	}

	if cfResp.Success {
		return true, nil
	}

	// Check if it's a "not found" vs some other error.
	for _, e := range cfResp.Errors {
		if e.Code == 8000007 { // Cloudflare's "project not found" code
			return false, nil
		}
	}
	return false, fmt.Errorf("Cloudflare API error: %s", formatErrors(cfResp.Errors))
}

// createProject creates a new Cloudflare Pages project.
func (c *CloudflareDeployer) createProject() error {
	url := fmt.Sprintf("%s/accounts/%s/pages/projects",
		cloudflareAPIBase, c.config.AccountID)

	body := fmt.Sprintf(`{"name":%q,"production_branch":"main"}`, c.config.Project)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

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
		return fmt.Errorf("failed to create project: %s", formatErrors(cfResp.Errors))
	}

	fmt.Printf("Created Cloudflare Pages project %q\n", c.config.Project)
	return nil
}

// Deploy uploads the output directory to Cloudflare Pages as a new deployment.
func (c *CloudflareDeployer) Deploy() error {
	if err := c.validate(); err != nil {
		return err
	}

	c.progress(0, "Checking project...")
	if err := c.ensureProject(); err != nil {
		return err
	}

	c.progress(2, "Collecting files...")

	files, err := c.collectFiles()
	if err != nil {
		return fmt.Errorf("collecting files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found in output directory: %s", c.output)
	}

	c.progress(5, fmt.Sprintf("Preparing %d files...", len(files)))

	// Build multipart body to a temp file so we can retry and get clean errors.
	tmpFile, err := os.CreateTemp("", "purtypics-deploy-*.bin")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := multipart.NewWriter(tmpFile)

	// Add branch field first if configured.
	if c.config.Branch != "" {
		if err := writer.WriteField("branch", c.config.Branch); err != nil {
			return fmt.Errorf("writing branch field: %w", err)
		}
	}

	for i, relPath := range files {
		absPath := filepath.Join(c.output, relPath)
		f, err := os.Open(absPath)
		if err != nil {
			return fmt.Errorf("opening %s: %w", relPath, err)
		}

		// Cloudflare expects the part name to be the path with leading slash.
		partName := "/" + relPath
		part, err := writer.CreateFormFile(partName, filepath.Base(relPath))
		if err != nil {
			f.Close()
			return fmt.Errorf("creating form part for %s: %w", relPath, err)
		}

		if _, err := io.Copy(part, f); err != nil {
			f.Close()
			return fmt.Errorf("writing %s: %w", relPath, err)
		}
		f.Close()

		pct := 5 + (45 * (i + 1) / len(files))
		c.progress(pct, fmt.Sprintf("Preparing %d/%d files...", i+1, len(files)))
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("finalizing multipart body: %w", err)
	}

	// Seek back to start for reading.
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("seeking temp file: %w", err)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		return fmt.Errorf("stat temp file: %w", err)
	}

	c.progress(55, fmt.Sprintf("Uploading %d files to Cloudflare...", len(files)))

	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s/deployments",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	req, err := http.NewRequest(http.MethodPost, apiURL, tmpFile)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = stat.Size()

	c.progress(60, "Uploading to Cloudflare...")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("uploading to Cloudflare: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return fmt.Errorf("decoding response (status %d): %s", resp.StatusCode, string(body))
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
