package deploy

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeebo/blake3"
)

const cloudflareAPIBase = "https://api.cloudflare.com/client/v4"

const (
	maxBucketSize      = 40 * 1024 * 1024 // 40 MiB per upload bucket
	maxBucketFileCount = 2000
	maxAssetSize       = 25 * 1024 * 1024 // 25 MiB per file
)

// CloudflareDeployer handles deployments to Cloudflare Pages.
type CloudflareDeployer struct {
	config           *CloudflareConfig
	output           string
	token            string // CLOUDFLARE_API_TOKEN
	jwt              string // short-lived upload JWT, refreshed as needed
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

type uploadTokenResult struct {
	JWT string `json:"jwt"`
}

type deploymentResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// fileEntry tracks a file for upload to Pages.
type fileEntry struct {
	path        string // relative path with forward slashes
	hash        string // 32-char hex blake3 hash
	b64Content  string // base64-encoded file content
	size        int64  // original file size
	contentType string
}

// largeFileEntry tracks a file too large for Pages (>25 MiB).
type largeFileEntry struct {
	path        string // relative path with forward slashes
	absPath     string // absolute filesystem path
	size        int64
	contentType string
}

// uploadItem is the JSON structure for a single file in an upload bucket.
type uploadItem struct {
	Key      string         `json:"key"`
	Value    string         `json:"value"`
	Metadata uploadMetadata `json:"metadata"`
	Base64   bool           `json:"base64"`
}

type uploadMetadata struct {
	ContentType string `json:"contentType"`
}

// LogCloudflareCredentials prints which Cloudflare-related env vars are detected.
func LogCloudflareCredentials(r2Enabled bool) {
	cfToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	r2Access := os.Getenv("R2_ACCESS_KEY_ID")
	r2Secret := os.Getenv("R2_SECRET_ACCESS_KEY")

	fmt.Println("Credential check:")
	if cfToken != "" {
		fmt.Printf("  CLOUDFLARE_API_TOKEN: set (%d chars)\n", len(cfToken))
	} else {
		fmt.Println("  CLOUDFLARE_API_TOKEN: NOT SET")
	}

	if r2Enabled {
		if r2Access != "" {
			fmt.Printf("  R2_ACCESS_KEY_ID:     set (%d chars)\n", len(r2Access))
		} else {
			fmt.Println("  R2_ACCESS_KEY_ID:     NOT SET")
		}
		if r2Secret != "" {
			fmt.Printf("  R2_SECRET_ACCESS_KEY: set (%d chars)\n", len(r2Secret))
		} else {
			fmt.Println("  R2_SECRET_ACCESS_KEY: NOT SET")
		}
	}
}

// NewCloudflareDeployer creates a new Cloudflare Pages deployer.
func NewCloudflareDeployer(config *CloudflareConfig, outputPath string) (*CloudflareDeployer, error) {
	r2Enabled := config.R2 != nil && config.R2.Enabled
	LogCloudflareCredentials(r2Enabled)

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

// TestConnection verifies the API token and project existence.
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

func (c *CloudflareDeployer) projectExists() (bool, error) {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	body, statusCode, err := c.apiGet(apiURL)
	if err != nil {
		return false, err
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return false, fmt.Errorf("decoding response (HTTP %d): %s", statusCode, string(body))
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

func (c *CloudflareDeployer) createProject() error {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects",
		cloudflareAPIBase, c.config.AccountID)

	payload := fmt.Sprintf(`{"name":%q,"production_branch":"main"}`, c.config.Project)

	body, err := c.apiPost(apiURL, "application/json", strings.NewReader(payload), c.token)
	if err != nil {
		return err
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return fmt.Errorf("decoding response: %s", string(body))
	}
	if !cfResp.Success {
		return fmt.Errorf("failed to create project: %s", formatErrors(cfResp.Errors))
	}

	fmt.Printf("Created Cloudflare Pages project %q\n", c.config.Project)
	return nil
}

// Deploy uploads the output directory to Cloudflare Pages.
//
// Flow:
//  1. Hash files, separating small (≤25 MiB) from large (>25 MiB)
//  2. If R2 enabled: upload large files to R2, rewrite HTML references
//  3. Upload small files to Pages (upload-token → check-missing → upload → upsert-hashes)
//  4. Create deployment with manifest
func (c *CloudflareDeployer) Deploy() error {
	if err := c.validate(); err != nil {
		return err
	}

	c.progress(0, "Checking project...")
	if err := c.ensureProject(); err != nil {
		return err
	}

	c.progress(2, "Scanning files...")
	entries, largeFiles, err := c.collectFiles()
	if err != nil {
		return fmt.Errorf("scanning files: %w", err)
	}
	if len(entries) == 0 && len(largeFiles) == 0 {
		return fmt.Errorf("no files found in output directory: %s", c.output)
	}

	// Handle large files: upload to R2 or warn.
	if len(largeFiles) > 0 {
		if err := c.handleLargeFiles(largeFiles, entries); err != nil {
			return err
		}
	}

	if len(entries) == 0 {
		return fmt.Errorf("no files to deploy to Pages (all files exceed 25 MiB limit)")
	}

	// Build manifest and collect all hashes.
	manifest := make(map[string]string, len(entries))
	allHashes := make([]string, 0, len(entries))
	for _, e := range entries {
		manifest["/"+e.path] = e.hash
		allHashes = append(allHashes, e.hash)
	}

	c.progress(30, "Getting upload token...")
	if err := c.refreshJWT(); err != nil {
		return fmt.Errorf("getting upload token: %w", err)
	}

	c.progress(33, "Checking which files need uploading...")
	missing, err := c.checkMissing(allHashes)
	if err != nil {
		return fmt.Errorf("checking missing files: %w", err)
	}

	// Filter to only files that need uploading.
	missingSet := make(map[string]bool, len(missing))
	for _, h := range missing {
		missingSet[h] = true
	}
	var toUpload []fileEntry
	for _, e := range entries {
		if missingSet[e.hash] {
			toUpload = append(toUpload, e)
		}
	}

	c.progress(36, fmt.Sprintf("Uploading %d/%d files to Pages...", len(toUpload), len(entries)))

	if len(toUpload) > 0 {
		if err := c.uploadBuckets(toUpload, len(entries)); err != nil {
			return fmt.Errorf("uploading files: %w", err)
		}
	}

	c.progress(88, "Registering files...")
	if err := c.upsertHashes(allHashes); err != nil {
		return fmt.Errorf("registering hashes: %w", err)
	}

	c.progress(92, "Creating deployment...")
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

// handleLargeFiles either uploads large files to R2 and rewrites HTML, or warns.
func (c *CloudflareDeployer) handleLargeFiles(largeFiles []largeFileEntry, entries []fileEntry) error {
	r2Enabled := c.config.R2 != nil && c.config.R2.Enabled

	if !r2Enabled {
		var totalSize int64
		for _, f := range largeFiles {
			totalSize += f.size
			fmt.Printf("  Skipping %s (%.1f MiB)\n", f.path, float64(f.size)/(1024*1024))
		}
		fmt.Printf("WARNING: %d file(s) (%.1f MiB total) exceed the 25 MiB Pages limit and will not be deployed.\n",
			len(largeFiles), float64(totalSize)/(1024*1024))
		fmt.Println("Enable R2 in deployment settings to upload large files to Cloudflare R2.")
		return nil
	}

	ctx := context.Background()

	c.progress(5, "Connecting to R2...")
	r2, err := NewR2Client(c.config.AccountID)
	if err != nil {
		return fmt.Errorf("R2 setup: %w", err)
	}

	bucket := c.config.R2.Bucket
	if bucket == "" {
		bucket = c.config.Project + "-assets"
	}

	c.progress(7, "Ensuring R2 bucket...")
	if err := r2.EnsureBucket(ctx, bucket); err != nil {
		return fmt.Errorf("R2 bucket: %w", err)
	}

	if err := r2.EnablePublicAccess(ctx, bucket); err != nil {
		return fmt.Errorf("R2 public access: %w", err)
	}

	publicBaseURL, err := r2.GetPublicURL(ctx, bucket, c.config.R2.CustomDomain)
	if err != nil {
		return fmt.Errorf("R2 public URL: %w", err)
	}
	publicBaseURL = strings.TrimRight(publicBaseURL, "/")

	// Upload large files to R2.
	replacements := make(map[string]string) // relative path -> R2 URL
	for i, f := range largeFiles {
		c.progress(10+15*i/len(largeFiles),
			fmt.Sprintf("Uploading to R2: %s (%.1f MiB)...", f.path, float64(f.size)/(1024*1024)))

		key := f.path // use same relative path as key
		if err := r2.UploadFile(ctx, bucket, key, f.absPath); err != nil {
			return fmt.Errorf("R2 upload %s: %w", f.path, err)
		}

		r2URL := publicBaseURL + "/" + key
		replacements[f.path] = r2URL
		fmt.Printf("  Uploaded %s -> %s\n", f.path, r2URL)
	}

	// Rewrite HTML files to reference R2 URLs instead of local paths.
	if len(replacements) > 0 {
		c.progress(25, "Rewriting HTML references...")
		c.rewriteHTMLEntries(entries, replacements)
	}

	return nil
}

// rewriteHTMLEntries replaces local file paths with R2 URLs in HTML file entries.
func (c *CloudflareDeployer) rewriteHTMLEntries(entries []fileEntry, replacements map[string]string) {
	for i, entry := range entries {
		if !strings.HasSuffix(entry.path, ".html") {
			continue
		}

		// Decode the base64 content.
		content, err := base64.StdEncoding.DecodeString(entry.b64Content)
		if err != nil {
			continue
		}

		html := string(content)
		modified := false

		for relPath, r2URL := range replacements {
			// Templates use paths like "../static/videos/album/file.mov"
			// or "/static/videos/album/file.mov". Replace all variants.
			// Order matters: longer prefixes first to avoid re-matching
			// inside already-replaced R2 URLs (e.g. "/path" matching
			// inside "https://...r2.dev/path").
			variants := []string{
				"../" + relPath,         // from album pages: ../static/videos/...
				"./" + relPath,          // from index: ./static/videos/...
				"/" + relPath,           // absolute: /static/videos/...
			}
			for _, v := range variants {
				if strings.Contains(html, v) {
					html = strings.ReplaceAll(html, v, r2URL)
					modified = true
					break // avoid shorter variants matching inside the R2 URL
				}
			}
			// Also handle bare (unquoted) paths in attributes.
			bare := "\"" + relPath + "\""
			if strings.Contains(html, bare) {
				html = strings.ReplaceAll(html, bare, "\""+r2URL+"\"")
				modified = true
			}
		}

		if modified {
			newB64 := base64.StdEncoding.EncodeToString([]byte(html))

			// Recompute blake3 hash.
			ext := filepath.Ext(entry.path)
			if len(ext) > 0 {
				ext = ext[1:]
			}
			h := blake3.Sum256([]byte(newB64 + ext))
			newHash := hex.EncodeToString(h[:])[:32]

			entries[i].b64Content = newB64
			entries[i].hash = newHash
			entries[i].size = int64(len(html))
		}
	}
}

// collectFiles walks the output directory, hashing small files for Pages
// and collecting large files separately.
func (c *CloudflareDeployer) collectFiles() ([]fileEntry, []largeFileEntry, error) {
	var entries []fileEntry
	var largeFiles []largeFileEntry

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

		ct := mime.TypeByExtension(filepath.Ext(path))
		if ct == "" {
			ct = "application/octet-stream"
		}

		// Large files go to the separate list.
		if info.Size() > maxAssetSize {
			largeFiles = append(largeFiles, largeFileEntry{
				path:        rel,
				absPath:     path,
				size:        info.Size(),
				contentType: ct,
			})
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		b64 := base64.StdEncoding.EncodeToString(data)

		// Hash: blake3(base64Content + extension), first 32 hex chars.
		ext := filepath.Ext(path)
		if len(ext) > 0 {
			ext = ext[1:] // strip leading dot
		}
		h := blake3.Sum256([]byte(b64 + ext))
		hash := hex.EncodeToString(h[:])[:32]

		entries = append(entries, fileEntry{
			path:        rel,
			hash:        hash,
			b64Content:  b64,
			size:        info.Size(),
			contentType: ct,
		})
		return nil
	})
	return entries, largeFiles, err
}

// getUploadToken retrieves a JWT for uploading and managing assets.
// refreshJWT fetches a fresh upload JWT and stores it on the deployer.
func (c *CloudflareDeployer) refreshJWT() error {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s/upload-token",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	body, statusCode, err := c.apiGet(apiURL)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("get upload token failed (HTTP %d): %s", statusCode, string(body))
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return fmt.Errorf("decoding response: %s", string(body))
	}
	if !cfResp.Success {
		return fmt.Errorf("get upload token failed: %s", formatErrors(cfResp.Errors))
	}

	var result uploadTokenResult
	if err := json.Unmarshal(cfResp.Result, &result); err != nil {
		return fmt.Errorf("parsing upload token: %w", err)
	}
	if result.JWT == "" {
		return fmt.Errorf("empty upload token returned")
	}

	c.jwt = result.JWT
	return nil
}

// jwtPost performs an API POST using the upload JWT.
// On 403 (expired JWT), it refreshes the token and retries once.
func (c *CloudflareDeployer) jwtPost(url, contentType string, payload []byte) ([]byte, error) {
	body, err := c.apiPost(url, contentType, strings.NewReader(string(payload)), c.jwt)
	if err != nil && strings.Contains(err.Error(), "HTTP 403") {
		fmt.Println("Upload token expired, refreshing...")
		if refreshErr := c.refreshJWT(); refreshErr != nil {
			return nil, fmt.Errorf("refreshing upload token: %w", refreshErr)
		}
		return c.apiPost(url, contentType, strings.NewReader(string(payload)), c.jwt)
	}
	return body, err
}

// checkMissing asks Cloudflare which file hashes it doesn't already have.
func (c *CloudflareDeployer) checkMissing(hashes []string) ([]string, error) {
	apiURL := fmt.Sprintf("%s/pages/assets/check-missing", cloudflareAPIBase)

	payload, _ := json.Marshal(map[string]interface{}{
		"hashes": hashes,
	})

	body, err := c.jwtPost(apiURL, "application/json", payload)
	if err != nil {
		return nil, fmt.Errorf("check-missing request: %w", err)
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("decoding response: %s", string(body))
	}
	if !cfResp.Success {
		return nil, fmt.Errorf("check-missing failed: %s", formatErrors(cfResp.Errors))
	}

	var missing []string
	if err := json.Unmarshal(cfResp.Result, &missing); err != nil {
		return nil, fmt.Errorf("parsing missing hashes: %w", err)
	}

	return missing, nil
}

// uploadBuckets uploads files in size-limited buckets.
func (c *CloudflareDeployer) uploadBuckets(files []fileEntry, totalFiles int) error {
	uploaded := 0
	batchStart := 0

	for batchStart < len(files) {
		var bucket []fileEntry
		var bucketSize int64

		for i := batchStart; i < len(files); i++ {
			entrySize := int64(len(files[i].b64Content)) + 256 // JSON overhead estimate
			if len(bucket) > 0 && (len(bucket) >= maxBucketFileCount || bucketSize+entrySize > maxBucketSize) {
				break
			}
			bucket = append(bucket, files[i])
			bucketSize += entrySize
		}

		if err := c.uploadBucket(bucket); err != nil {
			return err
		}

		uploaded += len(bucket)
		pct := 36 + (52 * uploaded / totalFiles)
		c.progress(pct, fmt.Sprintf("Uploaded %d/%d files...", uploaded, len(files)))

		batchStart += len(bucket)
	}

	return nil
}

// uploadBucket uploads a single bucket of files as a JSON array.
// Retries up to 3 times on network errors with exponential backoff.
func (c *CloudflareDeployer) uploadBucket(bucket []fileEntry) error {
	items := make([]uploadItem, 0, len(bucket))
	for _, entry := range bucket {
		items = append(items, uploadItem{
			Key:      entry.hash,
			Value:    entry.b64Content,
			Metadata: uploadMetadata{ContentType: entry.contentType},
			Base64:   true,
		})
	}

	payload, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("marshaling upload bucket: %w", err)
	}

	apiURL := fmt.Sprintf("%s/pages/assets/upload", cloudflareAPIBase)

	const maxRetries = 3
	var lastErr error
	for attempt := range maxRetries {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * 5 * time.Second
			fmt.Printf("Retrying upload (attempt %d/%d) after %v...\n", attempt+1, maxRetries, delay)
			time.Sleep(delay)
		}

		body, err := c.jwtPost(apiURL, "application/json", payload)
		if err != nil {
			lastErr = fmt.Errorf("upload bucket: %w", err)
			continue
		}

		var cfResp cloudflareResponse
		if err := json.Unmarshal(body, &cfResp); err != nil {
			return fmt.Errorf("decoding response: %s", string(body))
		}
		if !cfResp.Success {
			return fmt.Errorf("upload failed: %s", formatErrors(cfResp.Errors))
		}

		return nil
	}

	return lastErr
}

// upsertHashes registers all file hashes with Cloudflare.
func (c *CloudflareDeployer) upsertHashes(hashes []string) error {
	payload, _ := json.Marshal(map[string]interface{}{
		"hashes": hashes,
	})

	apiURL := fmt.Sprintf("%s/pages/assets/upsert-hashes", cloudflareAPIBase)

	body, err := c.jwtPost(apiURL, "application/json", payload)
	if err != nil {
		return fmt.Errorf("upsert-hashes request: %w", err)
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return fmt.Errorf("decoding response: %s", string(body))
	}
	if !cfResp.Success {
		return fmt.Errorf("upsert-hashes failed: %s", formatErrors(cfResp.Errors))
	}

	return nil
}

// createDeployment creates a deployment with the file manifest.
func (c *CloudflareDeployer) createDeployment(manifest map[string]string) (string, error) {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s/deployments",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("marshaling manifest: %w", err)
	}

	// Build multipart form body.
	boundary := "----purtypics-deploy"
	var sb strings.Builder
	sb.WriteString("--" + boundary + "\r\n")
	sb.WriteString("Content-Disposition: form-data; name=\"manifest\"\r\n\r\n")
	sb.Write(manifestJSON)
	sb.WriteString("\r\n")

	if c.config.Branch != "" {
		sb.WriteString("--" + boundary + "\r\n")
		sb.WriteString("Content-Disposition: form-data; name=\"branch\"\r\n\r\n")
		sb.WriteString(c.config.Branch)
		sb.WriteString("\r\n")
	}

	sb.WriteString("--" + boundary + "--\r\n")

	body, err := c.apiPost(apiURL, "multipart/form-data; boundary="+boundary,
		strings.NewReader(sb.String()), c.token)
	if err != nil {
		return "", fmt.Errorf("creating deployment: %w", err)
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return "", fmt.Errorf("decoding response: %s", string(body))
	}
	if !cfResp.Success {
		return "", fmt.Errorf("create deployment failed: %s", formatErrors(cfResp.Errors))
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

// apiGet performs a GET request with the API token.
func (c *CloudflareDeployer) apiGet(url string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("connecting to Cloudflare API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return body, resp.StatusCode, nil
}

// apiPost performs a POST request with the given auth token.
func (c *CloudflareDeployer) apiPost(url, contentType string, payload io.Reader, authToken string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connecting to Cloudflare API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func formatErrors(errs []cloudflareError) string {
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Message
	}
	return strings.Join(msgs, "; ")
}
