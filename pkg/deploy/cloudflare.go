package deploy

import (
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

type uploadTokenResult struct {
	JWT string `json:"jwt"`
}

type deploymentResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// fileEntry tracks a file for upload.
type fileEntry struct {
	path        string // relative path with forward slashes
	hash        string // 32-char hex blake3 hash
	b64Content  string // base64-encoded file content
	size        int64  // original file size
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

// NewCloudflareDeployer creates a new Cloudflare Pages deployer.
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
// Follows the same flow as Wrangler:
//  1. Hash all files (blake3 of base64(content) + extension, truncated to 32 hex chars)
//  2. Get upload JWT
//  3. Check which files Cloudflare already has
//  4. Upload missing files in buckets
//  5. Register all hashes
//  6. Create deployment with manifest
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

	// Build manifest and collect all hashes.
	manifest := make(map[string]string, len(entries))
	allHashes := make([]string, 0, len(entries))
	for _, e := range entries {
		manifest["/"+e.path] = e.hash
		allHashes = append(allHashes, e.hash)
	}

	c.progress(5, "Getting upload token...")
	jwt, err := c.getUploadToken()
	if err != nil {
		return fmt.Errorf("getting upload token: %w", err)
	}

	c.progress(8, "Checking which files need uploading...")
	missing, err := c.checkMissing(jwt, allHashes)
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

	c.progress(12, fmt.Sprintf("Uploading %d/%d files...", len(toUpload), len(entries)))

	if len(toUpload) > 0 {
		if err := c.uploadBuckets(jwt, toUpload, len(entries)); err != nil {
			return fmt.Errorf("uploading files: %w", err)
		}
	}

	c.progress(82, "Registering files...")
	if err := c.upsertHashes(jwt, allHashes); err != nil {
		return fmt.Errorf("registering hashes: %w", err)
	}

	c.progress(88, "Creating deployment...")
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

// hashFiles walks the output directory, reads each file, and computes
// blake3 hashes matching Wrangler's algorithm: blake3(base64(content) + ext)[:32].
func (c *CloudflareDeployer) hashFiles() ([]fileEntry, error) {
	var entries []fileEntry
	err := filepath.Walk(c.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Size() > maxAssetSize {
			return fmt.Errorf("file %s exceeds 25 MiB limit (%d bytes)", path, info.Size())
		}

		rel, err := filepath.Rel(c.output, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if rel == "gallery.yaml" || rel == "deploy.yaml" {
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

		ct := mime.TypeByExtension(filepath.Ext(path))
		if ct == "" {
			ct = "application/octet-stream"
		}

		entries = append(entries, fileEntry{
			path:        rel,
			hash:        hash,
			b64Content:  b64,
			size:        info.Size(),
			contentType: ct,
		})
		return nil
	})
	return entries, err
}

// getUploadToken retrieves a JWT for uploading and managing assets.
func (c *CloudflareDeployer) getUploadToken() (string, error) {
	apiURL := fmt.Sprintf("%s/accounts/%s/pages/projects/%s/upload-token",
		cloudflareAPIBase, c.config.AccountID, c.config.Project)

	body, statusCode, err := c.apiGet(apiURL)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("get upload token failed (HTTP %d): %s", statusCode, string(body))
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return "", fmt.Errorf("decoding response: %s", string(body))
	}
	if !cfResp.Success {
		return "", fmt.Errorf("get upload token failed: %s", formatErrors(cfResp.Errors))
	}

	var result uploadTokenResult
	if err := json.Unmarshal(cfResp.Result, &result); err != nil {
		return "", fmt.Errorf("parsing upload token: %w", err)
	}
	if result.JWT == "" {
		return "", fmt.Errorf("empty upload token returned")
	}

	return result.JWT, nil
}

// checkMissing asks Cloudflare which file hashes it doesn't already have.
func (c *CloudflareDeployer) checkMissing(jwt string, hashes []string) ([]string, error) {
	apiURL := fmt.Sprintf("%s/pages/assets/check-missing", cloudflareAPIBase)

	payload, _ := json.Marshal(map[string]interface{}{
		"hashes": hashes,
	})

	body, err := c.apiPost(apiURL, "application/json", strings.NewReader(string(payload)), jwt)
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
func (c *CloudflareDeployer) uploadBuckets(jwt string, files []fileEntry, totalFiles int) error {
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

		if err := c.uploadBucket(jwt, bucket); err != nil {
			return err
		}

		uploaded += len(bucket)
		pct := 12 + (70 * uploaded / totalFiles)
		c.progress(pct, fmt.Sprintf("Uploaded %d/%d files...", uploaded, len(files)))

		batchStart += len(bucket)
	}

	return nil
}

// uploadBucket uploads a single bucket of files as a JSON array.
func (c *CloudflareDeployer) uploadBucket(jwt string, bucket []fileEntry) error {
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

	body, err := c.apiPost(apiURL, "application/json", strings.NewReader(string(payload)), jwt)
	if err != nil {
		return fmt.Errorf("upload bucket: %w", err)
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

// upsertHashes registers all file hashes with Cloudflare.
func (c *CloudflareDeployer) upsertHashes(jwt string, hashes []string) error {
	payload, _ := json.Marshal(map[string]interface{}{
		"hashes": hashes,
	})

	apiURL := fmt.Sprintf("%s/pages/assets/upsert-hashes", cloudflareAPIBase)

	body, err := c.apiPost(apiURL, "application/json", strings.NewReader(string(payload)), jwt)
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
