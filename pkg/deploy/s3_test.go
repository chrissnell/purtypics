package deploy

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// mockS3 implements s3API for testing.
type mockS3 struct {
	objects   map[string]string // key -> etag (unquoted)
	uploaded  map[string][]byte // key -> body bytes
	deleted   []string
	headErr   error
	putErr    error
	deleteErr error
	listErr   error
}

func newMockS3() *mockS3 {
	return &mockS3{
		objects:  make(map[string]string),
		uploaded: make(map[string][]byte),
	}
}

func (m *mockS3) HeadBucket(_ context.Context, _ *s3.HeadBucketInput, _ ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, m.headErr
}

func (m *mockS3) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if m.putErr != nil {
		return nil, m.putErr
	}
	body, _ := io.ReadAll(input.Body)
	m.uploaded[*input.Key] = body
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3) DeleteObject(_ context.Context, input *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if m.deleteErr != nil {
		return nil, m.deleteErr
	}
	m.deleted = append(m.deleted, *input.Key)
	return &s3.DeleteObjectOutput{}, nil
}

func (m *mockS3) ListObjectsV2(_ context.Context, input *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var contents []s3types.Object
	for k, etag := range m.objects {
		quotedETag := "\"" + etag + "\""
		contents = append(contents, s3types.Object{
			Key:  aws.String(k),
			ETag: aws.String(quotedETag),
		})
	}
	return &s3.ListObjectsV2Output{
		Contents:    contents,
		IsTruncated: aws.Bool(false),
	}, nil
}

// helper to compute MD5 hex of a byte slice
func md5hex(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}

// helper to create a temp output dir with files
func setupOutputDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for path, content := range files {
		abs := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestCollectFiles(t *testing.T) {
	dir := setupOutputDir(t, map[string]string{
		"index.html":           "<html>hello</html>",
		"static/style.css":     "body{}",
		"album/photo.jpg":      "jpegdata",
		"gallery.yaml":         "title: test",  // should be excluded
		"deploy.yaml":          "s3: bucket: x", // should be excluded
	})

	d := NewS3Deployer(&S3Config{Bucket: "b", Region: "r"}, dir)
	files, err := d.collectFiles()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(files)
	expected := []string{"album/photo.jpg", "index.html", "static/style.css"}
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(files), files)
	}
	for i, f := range files {
		if f != expected[i] {
			t.Errorf("file[%d] = %q, want %q", i, f, expected[i])
		}
	}
}

func TestCollectFiles_Empty(t *testing.T) {
	dir := t.TempDir()
	d := NewS3Deployer(&S3Config{Bucket: "b", Region: "r"}, dir)
	files, err := d.collectFiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestContentType(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"photo.webp", "image/webp"},
		{"photo.avif", "image/avif"},
		{"font.woff2", "font/woff2"},
		{"page.html", "text/html"},
		{"data.json", "application/json"},
		{"unknown.qqq", "application/octet-stream"},
	}
	for _, tt := range tests {
		got := contentType(tt.path)
		// mime.TypeByExtension may include charset params, so just check prefix
		if tt.want == "application/octet-stream" {
			if got != tt.want {
				t.Errorf("contentType(%q) = %q, want %q", tt.path, got, tt.want)
			}
		} else if len(got) < len(tt.want) || got[:len(tt.want)] != tt.want {
			t.Errorf("contentType(%q) = %q, want prefix %q", tt.path, got, tt.want)
		}
	}
}

func TestValidate(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		config  *S3Config
		output  string
		wantErr string
	}{
		{"missing bucket", &S3Config{Region: "us-east-1"}, dir, "bucket name is required"},
		{"missing region", &S3Config{Bucket: "b"}, dir, "region is required"},
		{"missing output", &S3Config{Bucket: "b", Region: "r"}, "/nonexistent/path", "output directory not found"},
		{"valid", &S3Config{Bucket: "b", Region: "r"}, dir, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewS3Deployer(tt.config, tt.output)
			err := d.validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !contains(err.Error(), tt.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestMD5File(t *testing.T) {
	content := "hello world"
	dir := setupOutputDir(t, map[string]string{"test.txt": content})

	d := NewS3Deployer(&S3Config{Bucket: "b", Region: "r"}, dir)
	got, err := d.md5File("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := md5hex([]byte(content))
	if got != want {
		t.Errorf("md5File = %q, want %q", got, want)
	}
}

func TestDeploy_SkipsUnchanged(t *testing.T) {
	content := "hello"
	dir := setupOutputDir(t, map[string]string{"index.html": content})

	mock := newMockS3()
	mock.objects["index.html"] = md5hex([]byte(content)) // same hash

	d := NewS3Deployer(&S3Config{Bucket: "test", Region: "us-east-1"}, dir)
	d.client = mock

	if err := d.Deploy(); err != nil {
		t.Fatal(err)
	}

	if len(mock.uploaded) != 0 {
		t.Errorf("expected 0 uploads, got %d: %v", len(mock.uploaded), keys(mock.uploaded))
	}
	if len(mock.deleted) != 0 {
		t.Errorf("expected 0 deletes, got %d", len(mock.deleted))
	}
}

func TestDeploy_UploadsChanged(t *testing.T) {
	dir := setupOutputDir(t, map[string]string{
		"index.html": "new content",
		"style.css":  "body{}",
	})

	mock := newMockS3()
	mock.objects["index.html"] = "oldmd5hash" // different hash -> upload
	mock.objects["style.css"] = md5hex([]byte("body{}")) // same -> skip

	d := NewS3Deployer(&S3Config{Bucket: "test", Region: "us-east-1"}, dir)
	d.client = mock

	if err := d.Deploy(); err != nil {
		t.Fatal(err)
	}

	if _, ok := mock.uploaded["index.html"]; !ok {
		t.Error("expected index.html to be uploaded")
	}
	if _, ok := mock.uploaded["style.css"]; ok {
		t.Error("style.css should have been skipped (unchanged)")
	}
}

func TestDeploy_DeletesStale(t *testing.T) {
	dir := setupOutputDir(t, map[string]string{
		"index.html": "hello",
	})

	mock := newMockS3()
	mock.objects["index.html"] = md5hex([]byte("hello"))
	mock.objects["old-page.html"] = "abc123" // not local -> delete
	mock.objects["static/old.css"] = "def456" // not local -> delete

	d := NewS3Deployer(&S3Config{Bucket: "test", Region: "us-east-1"}, dir)
	d.client = mock

	if err := d.Deploy(); err != nil {
		t.Fatal(err)
	}

	sort.Strings(mock.deleted)
	if len(mock.deleted) != 2 {
		t.Fatalf("expected 2 deletes, got %d: %v", len(mock.deleted), mock.deleted)
	}
	if mock.deleted[0] != "old-page.html" || mock.deleted[1] != "static/old.css" {
		t.Errorf("unexpected deletes: %v", mock.deleted)
	}
}

func TestDeploy_FullSync(t *testing.T) {
	dir := setupOutputDir(t, map[string]string{
		"index.html":        "new index",
		"static/style.css":  "body{}",
		"album/photo.jpg":   "jpegdata",
	})

	mock := newMockS3()
	// index.html changed, style.css unchanged, photo.jpg is new, removed.html is stale
	mock.objects["index.html"] = "stale-hash"
	mock.objects["static/style.css"] = md5hex([]byte("body{}"))
	mock.objects["removed.html"] = "whatever"

	d := NewS3Deployer(&S3Config{Bucket: "test", Region: "us-east-1"}, dir)
	d.client = mock

	// Track progress messages
	var progressMsgs []string
	d.SetProgressCallback(func(pct int, msg string) {
		progressMsgs = append(progressMsgs, msg)
	})

	if err := d.Deploy(); err != nil {
		t.Fatal(err)
	}

	// Should upload index.html (changed) and album/photo.jpg (new)
	if len(mock.uploaded) != 2 {
		t.Errorf("expected 2 uploads, got %d: %v", len(mock.uploaded), keys(mock.uploaded))
	}
	if _, ok := mock.uploaded["index.html"]; !ok {
		t.Error("expected index.html upload")
	}
	if _, ok := mock.uploaded["album/photo.jpg"]; !ok {
		t.Error("expected album/photo.jpg upload")
	}

	// Should delete removed.html
	if len(mock.deleted) != 1 || mock.deleted[0] != "removed.html" {
		t.Errorf("expected delete of removed.html, got %v", mock.deleted)
	}

	// Should have reported progress
	if len(progressMsgs) == 0 {
		t.Error("expected progress messages")
	}
}

func TestDeploy_EmptyBucket(t *testing.T) {
	dir := setupOutputDir(t, map[string]string{
		"index.html": "hello",
	})

	mock := newMockS3()
	// Empty bucket, everything should upload

	d := NewS3Deployer(&S3Config{Bucket: "test", Region: "us-east-1"}, dir)
	d.client = mock

	if err := d.Deploy(); err != nil {
		t.Fatal(err)
	}

	if len(mock.uploaded) != 1 {
		t.Errorf("expected 1 upload, got %d", len(mock.uploaded))
	}
	if len(mock.deleted) != 0 {
		t.Errorf("expected 0 deletes, got %d", len(mock.deleted))
	}
}

func TestGetInfo(t *testing.T) {
	d := NewS3Deployer(&S3Config{Bucket: "my-bucket", Region: "eu-west-1"}, "/tmp")
	got := d.GetInfo()
	if !contains(got, "my-bucket") || !contains(got, "eu-west-1") {
		t.Errorf("GetInfo = %q, want bucket and region", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func keys(m map[string][]byte) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
