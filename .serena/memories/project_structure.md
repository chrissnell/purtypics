# Purtypics Project Structure

```
purtypics/
├── main.go                 # CLI entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksums
├── .gitignore             # Git ignore file
├── SESSION_SUMMARY.md     # Current session documentation
├── pkg/                   # Go packages (following standard Go project layout)
│   ├── gallery/           # Core gallery logic
│   │   ├── album.go       # Album/Photo structs and directory scanning
│   │   ├── generator.go   # Main gallery generation logic
│   │   ├── html.go        # HTML generation
│   │   └── templates.go   # CSS/JS/HTML templates
│   ├── image/
│   │   └── processor.go   # Image processing (thumbnails, etc.)
│   ├── video/
│   │   └── processor.go   # Video processing (thumbnails, previews)
│   ├── exif/
│   │   ├── extractor.go   # EXIF data extraction
│   │   └── parser.go      # EXIF parsing utilities
│   └── hugo/              # Hugo-related utilities (possibly legacy)
│       ├── theme.go
│       └── config.go
├── example/               # Sample photos/videos for testing
└── output/               # Generated gallery output (gitignored)
```

## Key Directories
- `pkg/`: All Go packages following standard layout
- `scratch/`: For temporary files and binaries (gitignored)
- `dist/`: For distribution builds (gitignored)
- `output/`: Default output directory for generated galleries

## Important Notes
- No test files currently exist in the codebase
- No Makefile or build scripts present
- Uses standard Go module system for dependencies
- Follows Go conventions for package organization under pkg/