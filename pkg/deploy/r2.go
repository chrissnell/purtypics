package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Client handles uploads to Cloudflare R2.
type R2Client struct {
	s3       *s3.Client
	apiToken string // CLOUDFLARE_API_TOKEN for R2 management APIs
	accountID string
}

// NewR2Client creates a client for Cloudflare R2 using S3-compatible API.
// Requires R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY env vars.
func NewR2Client(accountID string) (*R2Client, error) {
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY environment variables are required")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("loading R2 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &R2Client{
		s3:        client,
		apiToken:  apiToken,
		accountID: accountID,
	}, nil
}

// EnsureBucket creates the R2 bucket if it doesn't exist.
func (r *R2Client) EnsureBucket(ctx context.Context, bucket string) error {
	_, err := r.s3.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucket,
	})
	if err == nil {
		return nil
	}

	_, err = r.s3.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket,
	})
	if err != nil {
		return fmt.Errorf("creating R2 bucket %q: %w", bucket, err)
	}

	fmt.Printf("Created R2 bucket %q\n", bucket)
	return nil
}

// EnablePublicAccess enables the r2.dev managed public domain for the bucket.
func (r *R2Client) EnablePublicAccess(ctx context.Context, bucket string) error {
	if r.apiToken == "" {
		return fmt.Errorf("CLOUDFLARE_API_TOKEN required to enable R2 public access")
	}

	apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/r2/buckets/%s",
		r.accountID, bucket)

	// Enable public access on the bucket.
	payload := `{"public_access":{"enabled":true}}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, apiURL, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("enabling public access: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("enable public access failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetPublicURL returns the public base URL for the bucket.
// If customDomain is set, uses that. Otherwise fetches the r2.dev URL.
func (r *R2Client) GetPublicURL(ctx context.Context, bucket, customDomain string) (string, error) {
	if customDomain != "" {
		return "https://" + strings.TrimPrefix(customDomain, "https://"), nil
	}

	if r.apiToken == "" {
		return "", fmt.Errorf("CLOUDFLARE_API_TOKEN required to get R2 public URL")
	}

	// Get bucket details to find the r2.dev subdomain.
	apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/r2/buckets/%s",
		r.accountID, bucket)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getting bucket info: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("get bucket info failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response to find the public URL.
	var result struct {
		Result struct {
			PublicAccess struct {
				R2Dev struct {
					URL string `json:"url"`
				} `json:"r2_dev"`
			} `json:"public_access"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing bucket info: %w", err)
	}

	if result.Result.PublicAccess.R2Dev.URL != "" {
		return result.Result.PublicAccess.R2Dev.URL, nil
	}

	// Fallback: construct from account ID.
	return fmt.Sprintf("https://%s.r2.dev", bucket), nil
}

// UploadFile uploads a single file to R2. Skips if the file already exists with the same size.
func (r *R2Client) UploadFile(ctx context.Context, bucket, key, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat %s: %w", filePath, err)
	}

	// Check if already uploaded with same size.
	head, err := r.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err == nil && head.ContentLength != nil && *head.ContentLength == info.Size() {
		return nil // already uploaded
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening %s: %w", filePath, err)
	}
	defer f.Close()

	ct := mime.TypeByExtension(filepath.Ext(filePath))
	if ct == "" {
		ct = "application/octet-stream"
	}

	_, err = r.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        f,
		ContentType: &ct,
	})
	if err != nil {
		return fmt.Errorf("uploading to R2: %w", err)
	}

	return nil
}
