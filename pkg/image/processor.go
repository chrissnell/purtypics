package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/png"
)

// ThumbnailSizes defines standard thumbnail dimensions
type ThumbnailSizes struct {
	Small  int // for grid view
	Medium int // for larger grid
	Large  int // for detail view
	Full   int // max size for web display
}

// DefaultSizes returns recommended thumbnail sizes
func DefaultSizes() ThumbnailSizes {
	return ThumbnailSizes{
		Small:  400,
		Medium: 800,
		Large:  1600,
		Full:   2560, // good for up to 5K displays
	}
}

// Processor handles image operations
type Processor struct {
	sizes      ThumbnailSizes
	outputPath string
	quality    int
}

// NewProcessor creates a new image processor
func NewProcessor(outputPath string) *Processor {
	return &Processor{
		sizes:      DefaultSizes(),
		outputPath: outputPath,
		quality:    95, // High quality for gallery
	}
}

// ProcessImage generates all thumbnail sizes for an image
func (p *Processor) ProcessImage(sourcePath, albumID, photoID string) (map[string]string, error) {
	thumbnails := make(map[string]string)

	// Get source file info
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat source file: %w", err)
	}

	// Check for existing thumbnails first
	sizes := map[string]int{
		"small":  p.sizes.Small,
		"medium": p.sizes.Medium,
		"large":  p.sizes.Large,
		"full":   p.sizes.Full,
	}

	// Check if all thumbnails exist and are newer than source
	allCached := true
	for sizeName := range sizes {
		thumbDir := filepath.Join(p.outputPath, "static", "thumbs", albumID)
		ext := ".jpg"
		thumbPath := filepath.Join(thumbDir, fmt.Sprintf("%s_%s%s", photoID, sizeName, ext))
		relPath := filepath.Join("/static/thumbs", albumID, fmt.Sprintf("%s_%s%s", photoID, sizeName, ext))
		
		thumbInfo, err := os.Stat(thumbPath)
		if err != nil || thumbInfo.ModTime().Before(sourceInfo.ModTime()) {
			allCached = false
			break
		}
		thumbnails[sizeName] = relPath
	}

	// If all thumbnails are cached and up-to-date, return them
	if allCached {
		return thumbnails, nil
	}

	// Otherwise, regenerate all thumbnails
	thumbnails = make(map[string]string)

	// Decode image and auto-orient based on EXIF
	img, err := imaging.Open(sourcePath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Reuse sizes map from above

	for sizeName, maxDim := range sizes {
		// Skip if image is smaller than target
		if origWidth <= maxDim && origHeight <= maxDim && sizeName != "small" {
			continue
		}

		// Calculate new dimensions maintaining aspect ratio
		var newWidth, newHeight int
		if origWidth > origHeight {
			newWidth = maxDim
			newHeight = origHeight * maxDim / origWidth
		} else {
			newHeight = maxDim
			newWidth = origWidth * maxDim / origHeight
		}

		// Create output path
		thumbDir := filepath.Join(p.outputPath, "static", "thumbs", albumID)
		if err := os.MkdirAll(thumbDir, 0755); err != nil {
			return nil, err
		}

		ext := ".jpg" // always output JPEG for consistency
		thumbPath := filepath.Join(thumbDir, fmt.Sprintf("%s_%s%s", photoID, sizeName, ext))
		relPath := filepath.Join("/static/thumbs", albumID, fmt.Sprintf("%s_%s%s", photoID, sizeName, ext))

		// Resize with Lanczos filter for best quality
		resized := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
		
		// Save as JPEG with high quality
		out, err := os.Create(thumbPath)
		if err != nil {
			return nil, err
		}
		defer out.Close()

		if err := jpeg.Encode(out, resized, &jpeg.Options{Quality: p.quality}); err != nil {
			return nil, err
		}

		thumbnails[sizeName] = relPath
	}

	return thumbnails, nil
}

// GetImageDimensions returns width and height of an image
func GetImageDimensions(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = destination.ReadFrom(source)
	return err
}