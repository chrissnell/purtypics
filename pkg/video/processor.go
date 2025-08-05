package video

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Processor handles video operations

type Processor struct {
	outputPath string
}

// NewProcessor creates a new video processor
func NewProcessor(outputPath string) *Processor {
	return &Processor{
		outputPath: outputPath,
	}
}

// ExtractThumbnail extracts a frame from video as thumbnail
func (p *Processor) ExtractThumbnail(videoPath, albumID, photoID string) (string, error) {
	// Create output directory
	thumbDir := filepath.Join(p.outputPath, "static", "thumbs", albumID)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return "", err
	}

	// Output path for video thumbnail - strip extension from photoID if present
	basePhotoID := photoID
	if ext := filepath.Ext(photoID); ext != "" {
		basePhotoID = strings.TrimSuffix(photoID, ext)
	}
	thumbPath := filepath.Join(thumbDir, fmt.Sprintf("%s_poster.jpg", basePhotoID))
	relPath := path.Join("/static/thumbs", albumID, fmt.Sprintf("%s_poster.jpg", basePhotoID))

	// Check if thumbnail already exists
	if _, err := os.Stat(thumbPath); err == nil {
		return relPath, nil
	}

	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		// If ffmpeg not available, return a placeholder
		return "", fmt.Errorf("ffmpeg not found, cannot extract video thumbnail")
	}

	// Extract frame at 1 second (or 0 if video is shorter)
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		thumbPath,
	)

	if _, err := cmd.CombinedOutput(); err != nil {
		// Try at 0 seconds if 1 second fails
		cmd = exec.Command("ffmpeg",
			"-i", videoPath,
			"-ss", "00:00:00",
			"-vframes", "1",
			"-q:v", "2",
			"-y",
			thumbPath,
		)
		if output2, err2 := cmd.CombinedOutput(); err2 != nil {
			return "", fmt.Errorf("failed to extract thumbnail: %v, output: %s", err2, string(output2))
		}
	}

	return relPath, nil
}

// GetVideoDimensions returns width and height of a video
func GetVideoDimensions(videoPath string) (int, int, error) {
	// Check if ffprobe is available
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return 0, 0, fmt.Errorf("ffprobe not found")
	}

	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	// Parse output like "1920x1080"
	dims := strings.TrimSpace(string(output))
	parts := strings.Split(dims, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected dimension format: %s", dims)
	}

	var width, height int
	_, err = fmt.Sscanf(dims, "%dx%d", &width, &height)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

// CopyVideoToStatic copies the video file to the static directory
func (p *Processor) CopyVideoToStatic(videoPath, albumID, photoID string) (string, error) {
	// Create output directory
	videoDir := filepath.Join(p.outputPath, "static", "videos", albumID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return "", err
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(videoPath))
	
	var destPath, relPath string
	// Check if photoID already includes extension
	if strings.HasSuffix(strings.ToLower(photoID), ext) {
		destPath = filepath.Join(videoDir, photoID)
		relPath = path.Join("/static/videos", albumID, photoID)
	} else {
		destPath = filepath.Join(videoDir, fmt.Sprintf("%s%s", photoID, ext))
		relPath = path.Join("/static/videos", albumID, fmt.Sprintf("%s%s", photoID, ext))
	}

	// Check if already exists
	if _, err := os.Stat(destPath); err == nil {
		return relPath, nil
	}

	// Copy video file
	input, err := os.ReadFile(videoPath)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(destPath, input, 0644); err != nil {
		return "", err
	}

	return relPath, nil
}