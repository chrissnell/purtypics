package deploy

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Deployer handles deployments to Amazon S3.
type S3Deployer struct {
	config           *S3Config
	output           string
	progressCallback func(int, string)
}

// NewS3Deployer creates a new S3 deployer.
// AWS credentials are resolved by the SDK's default chain
// (env vars, shared credentials file, IAM role, etc.).
func NewS3Deployer(config *S3Config, outputPath string) *S3Deployer {
	return &S3Deployer{
		config: config,
		output: outputPath,
	}
}

// SetProgressCallback sets the progress callback function.
func (d *S3Deployer) SetProgressCallback(fn func(int, string)) {
	d.progressCallback = fn
}

// TestConnection verifies AWS credentials and bucket access.
func (d *S3Deployer) TestConnection() error {
	if err := d.validate(); err != nil {
		return err
	}

	ctx := context.Background()
	client, err := d.newClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &d.config.Bucket,
	})
	if err != nil {
		return fmt.Errorf("cannot access bucket %q: %w", d.config.Bucket, err)
	}

	return nil
}

// Deploy uploads the output directory to S3, then optionally invalidates CloudFront.
func (d *S3Deployer) Deploy() error {
	if err := d.validate(); err != nil {
		return err
	}

	ctx := context.Background()
	client, err := d.newClient(ctx)
	if err != nil {
		return err
	}

	d.progress(0, "Collecting files...")

	files, err := d.collectFiles()
	if err != nil {
		return fmt.Errorf("collecting files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found in output directory: %s", d.output)
	}

	d.progress(5, fmt.Sprintf("Uploading %d files to s3://%s...", len(files), d.config.Bucket))

	// Upload each file. Reserve 5-90% for uploads, 90-100% for CloudFront.
	for i, relPath := range files {
		if err := d.uploadFile(ctx, client, relPath); err != nil {
			return fmt.Errorf("uploading %s: %w", relPath, err)
		}
		pct := 5 + (85 * (i + 1) / len(files))
		d.progress(pct, fmt.Sprintf("Uploaded %d/%d files", i+1, len(files)))
	}

	// Invalidate CloudFront if configured.
	if d.config.CloudFrontID != "" {
		d.progress(92, "Invalidating CloudFront cache...")
		if err := d.invalidateCloudFront(ctx); err != nil {
			return fmt.Errorf("CloudFront invalidation: %w", err)
		}
	}

	d.progress(100, "Deployment complete")
	return nil
}

// GetInfo returns a human-readable description of the deployment target.
func (d *S3Deployer) GetInfo() string {
	return fmt.Sprintf("S3 bucket %q (%s)", d.config.Bucket, d.config.Region)
}

func (d *S3Deployer) validate() error {
	if d.config.Bucket == "" {
		return fmt.Errorf("S3 bucket name is required")
	}
	if d.config.Region == "" {
		return fmt.Errorf("S3 region is required")
	}
	if _, err := os.Stat(d.output); err != nil {
		return fmt.Errorf("output directory not found: %s", d.output)
	}
	return nil
}

func (d *S3Deployer) newClient(ctx context.Context) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(d.config.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}
	return s3.NewFromConfig(cfg), nil
}

func (d *S3Deployer) uploadFile(ctx context.Context, client *s3.Client, relPath string) error {
	absPath := filepath.Join(d.output, relPath)
	f, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use forward-slash key for S3.
	key := filepath.ToSlash(relPath)

	input := &s3.PutObjectInput{
		Bucket: &d.config.Bucket,
		Key:    &key,
		Body:   f,
	}

	// Set Content-Type from extension.
	if ct := contentType(relPath); ct != "" {
		input.ContentType = &ct
	}

	// Set Cache-Control if configured.
	if d.config.CacheControl != "" {
		input.CacheControl = &d.config.CacheControl
	}

	// Set storage class if configured.
	if d.config.StorageClass != "" {
		input.StorageClass = s3types.StorageClass(d.config.StorageClass)
	}

	// Set ACL if configured.
	if d.config.ACL != "" {
		input.ACL = s3types.ObjectCannedACL(d.config.ACL)
	}

	_, err = client.PutObject(ctx, input)
	return err
}

func (d *S3Deployer) invalidateCloudFront(ctx context.Context) error {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(d.config.Region),
	)
	if err != nil {
		return fmt.Errorf("loading AWS config: %w", err)
	}

	cfClient := cloudfront.NewFromConfig(cfg)

	ref := fmt.Sprintf("purtypics-%d", time.Now().Unix())
	_, err = cfClient.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: &d.config.CloudFrontID,
		InvalidationBatch: &cftypes.InvalidationBatch{
			CallerReference: &ref,
			Paths: &cftypes.Paths{
				Quantity: aws.Int32(1),
				Items:    []string{"/*"},
			},
		},
	})
	return err
}

func (d *S3Deployer) collectFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(d.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(d.output, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if rel == "gallery.yaml" || rel == "deploy.yaml" {
			return nil
		}

		files = append(files, rel)
		return nil
	})
	return files, err
}

func (d *S3Deployer) progress(pct int, msg string) {
	if d.progressCallback != nil {
		d.progressCallback(pct, msg)
	}
}

// contentType returns the MIME type for a file based on its extension.
func contentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	// Common types that mime.TypeByExtension may not know.
	switch ext {
	case ".webp":
		return "image/webp"
	case ".avif":
		return "image/avif"
	case ".woff2":
		return "font/woff2"
	}

	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	return "application/octet-stream"
}
