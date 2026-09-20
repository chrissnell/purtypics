package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/exif"
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

// sizeSpec pairs a size name with its target max dimension.
type sizeSpec struct {
	name   string
	maxDim int
}

func (p *Processor) sizeSpecs() [4]sizeSpec {
	return [4]sizeSpec{
		{"small", p.sizes.Small},
		{"medium", p.sizes.Medium},
		{"large", p.sizes.Large},
		{"full", p.sizes.Full},
	}
}

// ProcessImage generates all thumbnail sizes for an image and returns the
// (oriented) full-image dimensions. When every thumbnail is already present
// on disk and newer than the source, no decode is performed; dimensions are
// read from the source header instead.
func (p *Processor) ProcessImage(sourcePath, albumID, photoID string) (map[string]string, int, int, error) {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to stat source file: %w", err)
	}

	thumbDir := filepath.Join(p.outputPath, "static", "thumbs", albumID)
	specs := p.sizeSpecs()

	// Cache check: if every thumb exists and is newer than source, skip decode.
	thumbnails := make(map[string]string, len(specs))
	allCached := true
	for _, s := range specs {
		name := fmt.Sprintf("%s_%s.jpg", photoID, s.name)
		thumbPath := filepath.Join(thumbDir, name)
		info, err := os.Stat(thumbPath)
		if err != nil || info.ModTime().Before(sourceInfo.ModTime()) {
			allCached = false
			break
		}
		thumbnails[s.name] = path.Join("/static/thumbs", albumID, name)
	}

	if allCached {
		// Still need dimensions for templates. DecodeConfig reads headers only.
		w, h, err := orientedDimensions(sourcePath)
		if err != nil {
			return thumbnails, 0, 0, nil // dimensions are best-effort
		}
		return thumbnails, w, h, nil
	}

	// Regenerate: decode once, resize N times.
	img, err := imaging.Open(sourcePath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, 0, 0, err
	}
	bounds := img.Bounds()
	origWidth, origHeight := bounds.Dx(), bounds.Dy()

	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, 0, 0, err
	}

	thumbnails = make(map[string]string, len(specs))
	for _, s := range specs {
		// Don't upscale images beyond their native resolution; cap at native dimensions
		newWidth, newHeight := origWidth, origHeight
		if origWidth > s.maxDim || origHeight > s.maxDim {
			newWidth, newHeight = scaleToFit(origWidth, origHeight, s.maxDim)
		}
		name := fmt.Sprintf("%s_%s.jpg", photoID, s.name)
		thumbPath := filepath.Join(thumbDir, name)

		resized := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
		if err := saveJPEG(thumbPath, resized, p.quality); err != nil {
			return nil, 0, 0, err
		}
		thumbnails[s.name] = path.Join("/static/thumbs", albumID, name)
	}

	return thumbnails, origWidth, origHeight, nil
}

// scaleToFit computes dimensions that fit the source into a maxDim square
// while preserving aspect ratio.
func scaleToFit(w, h, maxDim int) (int, int) {
	if w > h {
		return maxDim, h * maxDim / w
	}
	return w * maxDim / h, maxDim
}

func saveJPEG(path string, img image.Image, quality int) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
}

// orientedDimensions returns width/height accounting for EXIF orientation.
// Uses DecodeConfig (header read only) plus a cheap orientation lookup.
func orientedDimensions(imgPath string) (int, int, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	// Swap dims for orientations that rotate 90/270 degrees.
	if o, err := exif.GetOrientation(imgPath); err == nil && o >= 5 && o <= 8 {
		return cfg.Height, cfg.Width, nil
	}
	return cfg.Width, cfg.Height, nil
}

// GetImageDimensions returns width and height of an image (raw, unoriented).
// Kept for backwards compatibility with callers that don't care about orientation.
func GetImageDimensions(imgPath string) (int, int, error) {
	file, err := os.Open(imgPath)
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
