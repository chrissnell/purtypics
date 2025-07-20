package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

// Load reads and parses metadata from a file
func Load(path string) (*GalleryMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty metadata if file doesn't exist
			return &GalleryMetadata{
				Albums: make(map[string]*AlbumMetadata),
				Photos: make(map[string]*PhotoMetadata),
			}, nil
		}
		return nil, fmt.Errorf("reading metadata file: %w", err)
	}

	meta := &GalleryMetadata{
		Albums: make(map[string]*AlbumMetadata),
		Photos: make(map[string]*PhotoMetadata),
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, meta)
	case ".json":
		err = json.Unmarshal(data, meta)
	default:
		return nil, fmt.Errorf("unsupported metadata format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("parsing metadata: %w", err)
	}

	return meta, nil
}

// Save writes metadata to a file
func Save(meta *GalleryMetadata, path string) error {
	var data []byte
	var err error

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(meta)
	case ".json":
		data, err = json.MarshalIndent(meta, "", "  ")
	default:
		return fmt.Errorf("unsupported metadata format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("marshaling metadata: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// GetAlbumMetadata returns metadata for a specific album
func (g *GalleryMetadata) GetAlbumMetadata(albumPath string) *AlbumMetadata {
	if g.Albums == nil {
		return nil
	}
	return g.Albums[albumPath]
}

// GetPhotoMetadata returns metadata for a specific photo
func (g *GalleryMetadata) GetPhotoMetadata(photoPath string) *PhotoMetadata {
	if g.Photos == nil {
		return nil
	}
	return g.Photos[photoPath]
}