package deploy

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

const cloudflareAPIBase = "https://api.cloudflare.com/client/v4"

// Max upload payload size per batch (25MB, well under Cloudflare's limit).
const maxBatchSize = 25 * 1024 * 1024

// CloudflareDeployer handles deployments to Cloudflare Pages.
type CloudflareDeployer struct {
	config           *CloudflareConfig
	output           string
	token            string
	progressCallback func(int, string)
}

// cloudflareResponse is the envelope for all Cloudflare API responses.
type cloudflareResponse struct {
	Success  bool                `json:"success"`
	Errors   []cloudflareError   `json:"errors"`
	Messages []cloudflareMessage `json:"messages"`
	Result   json.RawMessage     `json:"result"`
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

// uploadToken is returned by the check-missing endpoint.
type uploadToken struct {
	JWT string `json:"jwt"`
}

// fileEntry tracks a file's path, hash, and size for upload planning.
type fileEntry struct {
	path string // relative path with forward slashes
	hash string // hex-encoded SHA-256 content hash
	size int64
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
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
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

	for _, e := range cfResp.Errors {
		if e.Code == 8000007 {
			return false, nil
		}
	}
	return false, fmt.Errorf("Cloudflare API error: %s", formatErrors(cfResp.Errors))
}

// createProject creates a new Cloudflare Pages project.
func (c *CloudflareDeployer) createProject() error {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects",
		cloudflareAPIBase, c.config.AccountID)

	body := fmt.Sprintf(`{"name":%q,"production_branch":"main"}`, c.config.Project)

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(body))
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

// Deploy uploads the output directory to Cloudflare Pages using the Direct Upload API.
// This uses the batched upload flow: hash files, upload in chunks, create deployment.
func (c *CloudflareDeployer) Deploy() error {
	if err := c.validate(); err != nil {
		return err
	}

	c.progress(0, "Checking project...")
	if err := c.ensureProject(); err != nil {
		return err
	}

	c.progress(2, "Hashing files...")

	entries, err := c.hashFiles()
	if err != nil {
		return fmt.Errorf("hashing files: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("no files found in output directory: %s", c.output)
	}

	// Build the manifest: path -> hash.
	manifest := make(map[string]string, len(entries))
	hashes := make([]string, 0, len(entries))
	for _, e := range entries {
		manifest["/"+e.path] = e.hash
		hashes = append(hashes, e.hash)
	}

	c.progress(10, "Checking which files need uploading...")

	// Ask Cloudflare which files it already has.
	jwt, missing, err := c.checkMissing(hashes)
	if err != nil {
		return fmt.Errorf("checking missing files: %w", err)
	}

	// Build a set of missing hashes for quick lookup.
	missingSet := make(map[string]bool, len(missing))
	for _, h := range missing {
		missingSet[h] = true
	}

	// Filter entries to only those that need uploading.
	var toUpload []fileEntry
	for _, e := range entries {
		if missingSet[e.hash] {
			toUpload = append(toUpload, e)
		}
	}

	c.progress(15, fmt.Sprintf("Uploading %d/%d files...", len(toUpload), len(entries)))

	// Upload files in batches.
	if len(toUpload) > 0 {
		if err := c.uploadBatches(jwt, toUpload, len(entries)); err != nil {
			return err
		}
	}

	c.progress(85, "Creating deployment...")

	// Create the deployment with the manifest.
	deployURL, err := c.createDeployment(manifest)
	if err != nil {
		return fmt.Errorf("creating deployment: %w", err)
	}

	if deployURL != "" {
		c.progress(100, fmt.Sprintf("Deployed to %s", deployURL))
		fmt.Printf("Deployment URL: %s\n", deployURL)
	} else {
		c.progress(100, "Deployment complete")
	}

	return nil
}

// hashFiles walks the output directory and computes SHA-256 hashes.
func (c *CloudflareDeployer) hashFiles() ([]fileEntry, error) {
	var entries []fileEntry
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
		rel = filepath.ToSlash(rel)

		if rel == "gallery.yaml" || rel == "deploy.yaml" {
			return nil
		}

		h := sha256.New()
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		entries = append(entries, fileEntry{
			path: rel,
			hash: hex.EncodeToString(h.Sum(nil)),
			size: info.Size(),
		})
		return nil
	})
	return entries, err
}

// checkMissing sends file hashes to Cloudflare and returns a JWT for uploading
// and the list of hashes that Cloudflare doesn't already have.
func (c *CloudflareDeployer) checkMissing(hashes []string) (jwt string, missing []string, err error) {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/assets/check-missing",
		cloudflareAPIBase, c.config.AccountID)

	payload, _ := json.Marshal(map[string]interface{}{
		"hashes": hashes,
	})

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(payload))
	if err != nil {
		return "", nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("connecting to Cloudflare API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("check-missing failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return "", nil, fmt.Errorf("decoding response (status %d): %s", resp.StatusCode, string(body))
	}

	if !cfResp.Success {
		return "", nil, fmt.Errorf("check-missing failed: %s", formatErrors(cfResp.Errors))
	}

	// Parse the JWT from messages.
	for _, msg := range cfResp.Messages {
		if msg.Message != "" {
			jwt = msg.Message
			break
		}
	}
	if jwt == "" {
		return "", nil, fmt.Errorf("no upload JWT returned from check-missing")
	}

	// Parse missing hashes from result.
	if err := json.Unmarshal(cfResp.Result, &missing); err != nil {
		return "", nil, fmt.Errorf("parsing missing hashes: %w", err)
	}

	return jwt, missing, nil
}

// uploadBatches uploads files in size-limited batches using the JWT.
func (c *CloudflareDeployer) uploadBatches(jwt string, files []fileEntry, totalFiles int) error {
	uploaded := 0
	batchStart := 0

	for batchStart < len(files) {
		// Build a batch that fits within maxBatchSize.
		var batchSize int64
		batchEnd := batchStart
		for batchEnd < len(files) {
			if batchSize+files[batchEnd].size > maxBatchSize && batchEnd > batchStart {
				break
			}
			batchSize += files[batchEnd].size
			batchEnd++
		}

		batch := files[batchStart:batchEnd]
		if err := c.uploadBatch(jwt, batch); err != nil {
			return fmt.Errorf("uploading batch: %w", err)
		}

		uploaded += len(batch)
		pct := 15 + (70 * uploaded / totalFiles)
		c.progress(pct, fmt.Sprintf("Uploaded %d/%d files...", uploaded, len(files)))

		batchStart = batchEnd
	}

	return nil
}

// uploadBatch uploads a single batch of files as multipart form data.
func (c *CloudflareDeployer) uploadBatch(jwt string, batch []fileEntry) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, entry := range batch {
		absPath := filepath.Join(c.output, filepath.FromSlash(entry.path))
		f, err := os.Open(absPath)
		if err != nil {
			return fmt.Errorf("opening %s: %w", entry.path, err)
		}

		// Use the content hash as the part name (Cloudflare identifies files by hash).
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, entry.hash, filepath.Base(entry.path)))
		h.Set("Content-Type", "application/octet-stream")

		part, err := writer.CreatePart(h)
		if err != nil {
			f.Close()
			return fmt.Errorf("creating form part for %s: %w", entry.path, err)
		}

		if _, err := io.Copy(part, f); err != nil {
			f.Close()
			return fmt.Errorf("writing %s: %w", entry.path, err)
		}
		f.Close()
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("finalizing multipart: %w", err)
	}

	apiURL := fmt.Sprintf("%s/accounts/%s/pages/assets/upload",
		cloudflareAPIBase, c.config.AccountID)

	req, err := http.NewRequest(http.MethodPost, apiURL, &buf)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("uploading batch: %w", err)
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
		return fmt.Errorf("upload failed: %s", formatErrors(cfResp.Errors))
	}

	return nil
}

// createDeployment creates a deployment with the file manifest.
func (c *CloudflareDeployer) createDeployment(manifest map[string]string) (string, error) {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s/deployments",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add manifest as a JSON field.
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("marshaling manifest: %w", err)
	}
	if err := writer.WriteField("manifest", string(manifestJSON)); err != nil {
		return "", fmt.Errorf("writing manifest field: %w", err)
	}

	// Add branch if configured.
	if c.config.Branch != "" {
		if err := writer.WriteField("branch", c.config.Branch); err != nil {
			return "", fmt.Errorf("writing branch field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("finalizing multipart: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, &buf)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("creating deployment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return "", fmt.Errorf("decoding response (status %d): %s", resp.StatusCode, string(body))
	}

	if !cfResp.Success {
		return "", fmt.Errorf("deployment creation failed: %s", formatErrors(cfResp.Errors))
	}

	var result deploymentResult
	if err := json.Unmarshal(cfResp.Result, &result); err == nil {
		return result.URL, nil
	}
	return "", nil
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
