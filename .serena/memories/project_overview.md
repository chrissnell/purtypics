# Purtypics Project Overview

## Purpose
Purtypics is a static photo gallery generator written in Go that creates beautiful, responsive photo galleries with a masonry layout. It's inspired by Hugo but focused specifically on photo galleries with high-quality image processing.

## Tech Stack
- **Language**: Go 1.24
- **Dependencies**:
  - github.com/disintegration/imaging v1.6.2 (image processing)
  - github.com/rwcarlsen/goexif (EXIF data extraction)
  - golang.org/x/image (additional image utilities)
  - gopkg.in/yaml.v3 (YAML parsing)
- **External Requirements**: ffmpeg (for video support)

## Key Features
- High-quality thumbnail generation with Lanczos filtering
- EXIF metadata extraction and display
- Smart caching (only regenerates when source files change)
- Responsive masonry layout using desandro/masonry library
- Video support with thumbnail extraction and hover preview
- Lightbox for full-screen viewing
- Auto-rotation based on EXIF orientation
- Album sorting by newest first, photos by oldest first

## Usage
```bash
# Build the tool
go build -o purtypics

# Generate a gallery
./purtypics -source /path/to/photos -output /path/to/output -title "My Gallery"

# With base URL for deployment
./purtypics -source photos -output dist -title "My Gallery" -baseurl "https://example.com/gallery"
```

## Command-line Flags
- `-source`: Source directory containing photo albums (default: ".")
- `-output`: Output directory (defaults to source directory)
- `-baseurl`: Base URL for the site
- `-title`: Site title (default: "Photo Gallery")
- `-version`: Show version
- `-verbose`: Verbose output