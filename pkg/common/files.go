package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst, creating directories as needed.
func CopyFile(src, dst string) error {
	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer in.Close()

	// Create destination directory if needed
	if err := EnsureDirectory(filepath.Dir(dst)); err != nil {
		return err
	}

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close destination file: %w", cerr)
		}
	}()

	// Copy file contents
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Copy file permissions
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}
	if err := os.Chmod(dst, info.Mode()); err != nil {
		return fmt.Errorf("failed to set destination file permissions: %w", err)
	}

	return nil
}

// IsImageFile checks if a file has an image extension.
func IsImageFile(name string) bool {
	ext := filepath.Ext(name)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	case ".JPG", ".JPEG", ".PNG", ".GIF", ".BMP", ".TIFF", ".WEBP":
		return true
	}
	return false
}

// IsVideoFile checks if a file has a video extension.
func IsVideoFile(name string) bool {
	ext := filepath.Ext(name)
	switch ext {
	case ".mp4", ".mov", ".avi", ".webm", ".mkv", ".flv", ".wmv":
		return true
	case ".MP4", ".MOV", ".AVI", ".WEBM", ".MKV", ".FLV", ".WMV":
		return true
	}
	return false
}

// IsMediaFile checks if a file is an image or video.
func IsMediaFile(name string) bool {
	return IsImageFile(name) || IsVideoFile(name)
}