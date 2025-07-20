package editor

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjs/purtypics/pkg/exif"
	"github.com/disintegration/imaging"
)

// handleThumbnails dynamically generates thumbnails
func (s *Server) handleThumbnails(w http.ResponseWriter, r *http.Request) {
	// Extract path and size from URL
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/thumbs/"), "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid thumbnail path", http.StatusBadRequest)
		return
	}
	
	size := parts[0]
	imagePath := strings.Join(parts[1:], "/")
	fullPath := filepath.Join(s.SourcePath, imagePath)
	
	// Validate size
	var width int
	switch size {
	case "small":
		width = 300
	case "medium":
		width = 800
	case "large":
		width = 1600
	default:
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	
	// Check if source file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	
	// Open the image
	img, err := imaging.Open(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening image: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Apply EXIF orientation
	oriented, err := applyOrientation(fullPath, img)
	if err == nil {
		img = oriented
	}
	
	// Resize the image
	resized := imaging.Resize(img, width, 0, imaging.Lanczos)
	
	// Set content type based on file extension
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, resized, &jpeg.Options{Quality: 90}); err != nil {
			http.Error(w, "Error encoding image", http.StatusInternalServerError)
		}
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		if err := png.Encode(w, resized); err != nil {
			http.Error(w, "Error encoding image", http.StatusInternalServerError)
		}
	default:
		// Default to JPEG
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, resized, &jpeg.Options{Quality: 90}); err != nil {
			http.Error(w, "Error encoding image", http.StatusInternalServerError)
		}
	}
}

// applyOrientation applies EXIF orientation to an image
func applyOrientation(imagePath string, img image.Image) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return img, err
	}
	defer file.Close()

	orientation := 1
	exifData, err := exif.ExtractMetadata(imagePath)
	if err == nil && exifData != nil {
		orientation = exifData.Orientation
	}

	switch orientation {
	case 2: // Flip horizontal
		return imaging.FlipH(img), nil
	case 3: // Rotate 180
		return imaging.Rotate180(img), nil
	case 4: // Flip vertical
		return imaging.FlipV(img), nil
	case 5: // Transpose (flip horizontal and rotate 270)
		return imaging.Rotate270(imaging.FlipH(img)), nil
	case 6: // Rotate 90
		return imaging.Rotate270(img), nil
	case 7: // Transverse (flip horizontal and rotate 90)
		return imaging.Rotate90(imaging.FlipH(img)), nil
	case 8: // Rotate 270
		return imaging.Rotate90(img), nil
	default:
		return img, nil
	}
}