package deploy

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
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

// s3API abstracts the S3 operations we use, enabling test mocks.
type s3API interface {
	HeadBucket(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// S3Deployer handles deployments to Amazon S3.
type S3Deployer struct {
	config           *S3Config
	output           string
	client           s3API // nil until Deploy/TestConnection; injected in tests
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
	if err := d.ensureClient(ctx); err != nil {
		return err
	}

	_, err := d.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &d.config.Bucket,
	})
	if err != nil {
		return fmt.Errorf("cannot access bucket %q: %w", d.config.Bucket, err)
	}

	return nil
}

// Deploy syncs the output directory to S3:
//  1. Lists remote objects to build an ETag index.
//  2. Hashes local files (MD5) and compares against remote ETags.
//  3. Uploads only new or changed files.
//  4. Deletes remote objects that no longer exist locally.
//  5. Optionally invalidates CloudFront.
func (d *S3Deployer) Deploy() error {
	if err := d.validate(); err != nil {
		return err
	}

	ctx := context.Background()
	if err := d.ensureClient(ctx); err != nil {
		return err
	}

	d.progress(0, "Collecting local files...")

	localFiles, err := d.collectFiles()
	if err != nil {
		return fmt.Errorf("collecting files: %w", err)
	}
	if len(localFiles) == 0 {
		return fmt.Errorf("no files found in output directory: %s", d.output)
	}

	d.progress(5, "Listing remote objects...")

	remoteETags, err := d.listRemoteObjects(ctx)
	if err != nil {
		return fmt.Errorf("listing remote objects: %w", err)
	}

	d.progress(10, "Computing file hashes...")

	// Compute local MD5s and determine what needs uploading.
	localKeys := make(map[string]bool, len(localFiles))
	var toUpload []string
	for _, relPath := range localFiles {
		key := filepath.ToSlash(relPath)
		localKeys[key] = true

		localETag, err := d.md5File(relPath)
		if err != nil {
			return fmt.Errorf("hashing %s: %w", relPath, err)
		}

		if remoteETag, exists := remoteETags[key]; exists && remoteETag == localETag {
			continue // unchanged
		}
		toUpload = append(toUpload, relPath)
	}

	// Find remote objects to delete.
	var toDelete []string
	for key := range remoteETags {
		if !localKeys[key] {
			toDelete = append(toDelete, key)
		}
	}

	d.progress(15, fmt.Sprintf("Syncing: %d to upload, %d unchanged, %d to delete",
		len(toUpload), len(localFiles)-len(toUpload), len(toDelete)))

	// Upload new/changed files. Reserve 15-80% for uploads.
	for i, relPath := range toUpload {
		if err := d.uploadFile(ctx, relPath); err != nil {
			return fmt.Errorf("uploading %s: %w", relPath, err)
		}
		pct := 15 + (65 * (i + 1) / len(toUpload))
		d.progress(pct, fmt.Sprintf("Uploaded %d/%d changed files", i+1, len(toUpload)))
	}

	// Delete stale remote objects. Reserve 80-90%.
	for i, key := range toDelete {
		if err := d.deleteObject(ctx, key); err != nil {
			return fmt.Errorf("deleting %s: %w", key, err)
		}
		pct := 80 + (10 * (i + 1) / len(toDelete))
		d.progress(pct, fmt.Sprintf("Deleted %d/%d stale files", i+1, len(toDelete)))
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

func (d *S3Deployer) ensureClient(ctx context.Context) error {
	if d.client != nil {
		return nil
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(d.config.Region),
	)
	if err != nil {
		return fmt.Errorf("loading AWS config: %w", err)
	}
	d.client = s3.NewFromConfig(cfg)
	return nil
}

// listRemoteObjects returns a map of S3 key -> ETag for all objects in the bucket.
func (d *S3Deployer) listRemoteObjects(ctx context.Context) (map[string]string, error) {
	etags := make(map[string]string)
	var continuationToken *string

	for {
		out, err := d.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            &d.config.Bucket,
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}

		for _, obj := range out.Contents {
			if obj.Key != nil && obj.ETag != nil {
				// S3 ETags are quoted, strip quotes for comparison.
				etag := strings.Trim(*obj.ETag, "\"")
				etags[*obj.Key] = etag
			}
		}

		if !aws.ToBool(out.IsTruncated) {
			break
		}
		continuationToken = out.NextContinuationToken
	}

	return etags, nil
}

// md5File computes the hex-encoded MD5 of a local file (matches S3 ETag for non-multipart uploads).
func (d *S3Deployer) md5File(relPath string) (string, error) {
	absPath := filepath.Join(d.output, relPath)
	f, err := os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (d *S3Deployer) uploadFile(ctx context.Context, relPath string) error {
	absPath := filepath.Join(d.output, relPath)
	f, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer f.Close()

	key := filepath.ToSlash(relPath)

	input := &s3.PutObjectInput{
		Bucket: &d.config.Bucket,
		Key:    &key,
		Body:   f,
	}

	if ct := contentType(relPath); ct != "" {
		input.ContentType = &ct
	}
	if d.config.CacheControl != "" {
		input.CacheControl = &d.config.CacheControl
	}
	if d.config.StorageClass != "" {
		input.StorageClass = s3types.StorageClass(d.config.StorageClass)
	}
	if d.config.ACL != "" {
		input.ACL = s3types.ObjectCannedACL(d.config.ACL)
	}

	_, err = d.client.PutObject(ctx, input)
	return err
}

func (d *S3Deployer) deleteObject(ctx context.Context, key string) error {
	_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &d.config.Bucket,
		Key:    &key,
	})
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
