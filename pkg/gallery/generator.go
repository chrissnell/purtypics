package gallery

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/cjs/purtypics/pkg/exif"
	"github.com/cjs/purtypics/pkg/image"
	"github.com/cjs/purtypics/pkg/metadata"
	"github.com/cjs/purtypics/pkg/video"
)

// ProgressCallback is called to report generation progress
type ProgressCallback func(current, total int, message string)

// Generator creates the photo gallery
type Generator struct {
	SourcePath       string
	OutputPath       string
	SiteTitle        string
	BaseURL          string
	Verbose          bool
	MetadataPath     string
	metadata         *metadata.GalleryMetadata
	imageProcessor   *image.Processor
	videoProcessor   *video.Processor
	ProgressCallback ProgressCallback
}

// NewGenerator creates a new gallery generator
func NewGenerator(sourcePath, outputPath, siteTitle, baseURL string, verbose bool) *Generator {
	return &Generator{
		SourcePath:     sourcePath,
		OutputPath:     outputPath,
		SiteTitle:      siteTitle,
		BaseURL:        baseURL,
		Verbose:        verbose,
		imageProcessor: image.NewProcessor(outputPath),
		videoProcessor: video.NewProcessor(outputPath),
	}
}

// Generate processes all albums and creates the gallery
func (g *Generator) Generate() error {
	// Load metadata if path is provided
	if g.MetadataPath == "" {
		g.MetadataPath = filepath.Join(g.SourcePath, "gallery.yaml")
	}
	
	meta, err := metadata.Load(g.MetadataPath)
	if err != nil {
		return fmt.Errorf("loading metadata: %w", err)
	}
	g.metadata = meta
	
	// Apply gallery-level metadata
	if meta.Title != "" {
		g.SiteTitle = meta.Title
	}
	
	if g.Verbose {
		fmt.Printf("Scanning albums in %s\n", g.SourcePath)
	}

	// Find all albums
	albums, err := ScanDirectory(g.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(albums) == 0 {
		return fmt.Errorf("no albums found in %s", g.SourcePath)
	}

	fmt.Printf("Found %d albums\n", len(albums))

	// Report initial progress
	if g.ProgressCallback != nil {
		g.ProgressCallback(0, len(albums), "Starting album processing")
	}

	// Process each album and build filtered list
	filteredAlbums := make([]Album, 0, len(albums))
	for i := range albums {
		album := &albums[i]
		
		// Apply album metadata
		if albumMeta := g.metadata.GetAlbumMetadata(album.Path); albumMeta != nil {
			if albumMeta.Title != "" {
				album.Title = albumMeta.Title
			}
			if albumMeta.Description != "" {
				album.Description = albumMeta.Description
			}
			if albumMeta.Hidden {
				continue // Skip hidden albums
			}
			album.CoverPhoto = albumMeta.CoverPhoto
			if album.CoverPhoto != "" {
				fmt.Printf("  Using cover photo: %s\n", album.CoverPhoto)
			}
		}
		
		fmt.Printf("Processing %s\n", album.Title)
		
		// Report progress
		if g.ProgressCallback != nil {
			g.ProgressCallback(i+1, len(albums), fmt.Sprintf("Processing album: %s", album.Title))
		}
		
		if err := g.processAlbum(album); err != nil {
			log.Printf("Error processing album %s: %v", album.Title, err)
			continue
		}
		
		// Filter out hidden photos
		visiblePhotos := make([]Photo, 0, len(album.Photos))
		for _, photo := range album.Photos {
			if photo.Path != "" { // Non-empty path means it wasn't marked as hidden
				visiblePhotos = append(visiblePhotos, photo)
			}
		}
		album.Photos = visiblePhotos
		
		// Only include albums with visible photos
		if len(album.Photos) > 0 {
			// Sort photos by date after processing (oldest to newest)
			album.SortPhotosByDate()
			filteredAlbums = append(filteredAlbums, *album)
		}
	}
	albums = filteredAlbums

	// Sort albums by date (newest first) before generating HTML
	SortAlbumsByDate(albums)

	// Report HTML generation progress
	if g.ProgressCallback != nil {
		g.ProgressCallback(0, 1, "Generating HTML pages")
	}

	// Generate HTML site
	htmlGen := &HTMLGenerator{
		OutputPath:    g.OutputPath,
		SiteTitle:     g.SiteTitle,
		BaseURL:       g.BaseURL,
		ShowLocations: g.metadata.ShowLocations,
	}
	if err := htmlGen.Generate(albums); err != nil {
		return fmt.Errorf("failed to generate HTML site: %w", err)
	}

	// Report completion
	if g.ProgressCallback != nil {
		g.ProgressCallback(1, 1, "Gallery generation complete")
	}

	return nil
}

func (g *Generator) processAlbum(album *Album) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)
	processedCount := 0

	// Report photo processing start
	if g.ProgressCallback != nil {
		g.ProgressCallback(0, len(album.Photos), fmt.Sprintf("Processing photos in %s", album.Title))
	}

	// Limit concurrent processing
	sem := make(chan struct{}, 4)

	for i := range album.Photos {
		wg.Add(1)
		sem <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			photo := &album.Photos[idx]
			
			// Apply photo metadata
			photoPath := filepath.Join(album.Path, photo.Filename)
			if photoMeta := g.metadata.GetPhotoMetadata(photoPath); photoMeta != nil {
				if photoMeta.Title != "" {
					photo.Title = photoMeta.Title
				}
				if photoMeta.Description != "" {
					photo.Description = photoMeta.Description
				}
				if photoMeta.Hidden {
					// Mark photo for removal
					photo.Path = ""
					return
				}
			}

			if photo.IsVideo {
				// Handle video processing
				// Get video dimensions
				if width, height, err := video.GetVideoDimensions(photo.Path); err == nil {
					photo.Width = width
					photo.Height = height
				}

				// Extract video thumbnail
				if posterPath, err := g.videoProcessor.ExtractThumbnail(photo.Path, album.ID, photo.ID); err == nil {
					photo.Thumbnails = map[string]string{
						"poster": posterPath,
					}
					
					// Copy video to static directory
					if videoPath, err := g.videoProcessor.CopyVideoToStatic(photo.Path, album.ID, photo.ID); err == nil {
						photo.VideoPath = videoPath
					}
				} else {
					// If thumbnail extraction fails, skip this video
					mu.Lock()
					errors = append(errors, fmt.Errorf("video %s: %v", photo.Filename, err))
					mu.Unlock()
					return
				}
			} else {
				// Handle image processing
				// Get image dimensions
				if width, height, err := image.GetImageDimensions(photo.Path); err == nil {
					photo.Width = width
					photo.Height = height
				}

				// Extract EXIF data
				if exifData, err := exif.ExtractMetadata(photo.Path); err == nil {
					photo.EXIF = exifData
				}

				// Generate thumbnails
				thumbs, err := g.imageProcessor.ProcessImage(photo.Path, album.ID, photo.ID)
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
					return
				}
				photo.Thumbnails = thumbs
			}
			
			// Report progress
			mu.Lock()
			processedCount++
			if g.ProgressCallback != nil {
				g.ProgressCallback(processedCount, len(album.Photos), 
					fmt.Sprintf("Processing %s: %s", album.Title, photo.Filename))
			}
			mu.Unlock()
			
			if g.Verbose {
				fmt.Printf("  ✓ %s\n", photo.Filename)
			}
		}(i)
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors processing photos", len(errors))
	}

	return nil
}

