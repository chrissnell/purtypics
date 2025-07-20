# Purtypics

A fast, modern static photo gallery generator with built-in metadata editing.

## Features

- **Static Generation**: Creates fast, self-contained HTML galleries
- **Multi-Resolution**: Automatically generates multiple image sizes for optimal loading
- **Metadata Editor**: Built-in web interface for editing photo titles and descriptions
- **Video Support**: Handles videos with automatic thumbnail generation
- **EXIF Data**: Extracts and displays camera settings and location data
- **Responsive Design**: Beautiful masonry layout that works on all devices
- **Easy Deployment**: Deploy to any static host (rsync, S3, Cloudflare Pages)

## Quick Start

```bash
# Install purtypics
go install github.com/cjs/purtypics@latest

# Generate a gallery from your photos
purtypics generate /path/to/photos

# Edit metadata with the built-in editor
purtypics edit /path/to/photos

# Deploy to your server
purtypics deploy
```

## Installation

### From Source

```bash
git clone https://github.com/cjs/purtypics.git
cd purtypics
make build
```

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/cjs/purtypics/releases).

## Usage

### Basic Gallery Generation

The simplest way to create a gallery:

```bash
purtypics generate
```

This will:
1. Scan the current directory for photos and videos
2. Generate optimized thumbnails in multiple sizes
3. Create a static website in the `output/` directory

### Gallery Structure

Purtypics expects your photos to be organized in directories:

```
my-gallery/
├── vacation-2024/
│   ├── IMG_001.jpg
│   ├── IMG_002.jpg
│   └── VIDEO_001.mp4
├── family-reunion/
│   ├── DSC_001.jpg
│   └── DSC_002.jpg
└── gallery.yaml      # Optional metadata file
```

Each directory becomes an album in your gallery.

### Editing Metadata

Launch the web-based editor to add titles and descriptions:

```bash
purtypics edit
```

This opens a browser where you can:
- Add titles and descriptions to photos
- Mark favorites
- Hide photos from the gallery
- Set album cover photos
- Edit gallery title and settings

All changes are saved to `gallery.yaml`.

### Advanced Options

#### Custom Paths

```bash
# Specify source and output directories
purtypics generate -s /photos -o /website

# Edit with custom paths
purtypics edit -s /photos
```

#### Gallery Configuration

Create a `gallery.yaml` file to customize your gallery:

```yaml
title: "My Photo Collection"
description: "Family photos and adventures"

albums:
  vacation-2024:
    title: "Summer Vacation 2024"
    description: "Trip to the mountains"
    cover_photo: "IMG_001.jpg"
    
photos:
  vacation-2024/IMG_001.jpg:
    title: "Sunrise at the Summit"
    description: "Early morning hike was worth it!"
    favorite: true
```

### Deployment

#### rsync Deployment

Create a `deploy.yaml` file:

```yaml
rsync:
  host: user@yourserver.com
  path: /var/www/photos
  port: 22
```

Then deploy:

```bash
purtypics deploy
```

#### S3 Deployment (Coming Soon)

```yaml
s3:
  bucket: my-photo-gallery
  region: us-east-1
```

#### Cloudflare Pages (Coming Soon)

```yaml
cloudflare:
  project: my-gallery
  account_id: your-account-id
```

## Tips

### Performance

- **Original Files**: Keep your original photos in a separate backup. Purtypics generates optimized versions.
- **Large Collections**: For thousands of photos, organize into smaller albums for better performance.
- **Video Files**: Convert videos to MP4 format for best compatibility.

### Organization

- **Naming**: Use descriptive folder names - they become album titles
- **Dates**: Photos are sorted by EXIF date automatically
- **Hidden Files**: Files starting with `.` are ignored

### Quality Settings

Purtypics generates images in these sizes:
- **Small**: 300px wide (for thumbnails)
- **Medium**: 800px wide (for mobile devices)
- **Large**: 1600px wide (for desktop viewing)
- **Full**: Original resolution (for downloads)

## Troubleshooting

### Photos Not Appearing

- Check file extensions: `.jpg`, `.jpeg`, `.png`, `.webp`, `.heif`, `.heic`
- Ensure files aren't hidden (starting with `.`)
- Verify directory permissions

### Metadata Not Saving

- Check write permissions on `gallery.yaml`
- Ensure the editor can access the directory

### Deployment Fails

- Verify SSH access for rsync deployments
- Check credentials in `deploy.yaml`
- Ensure destination directory exists

## Privacy

Purtypics is designed for private photo collections:
- All processing happens locally on your machine
- No cloud services or external APIs are used
- EXIF location data can be excluded from generated files
- Password protection can be added at the web server level

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - See [LICENSE](LICENSE) for details.