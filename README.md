# Purtypics

A fast, modern static photo gallery generator with built-in metadata editing.

Check out a <a href="https://chrissnell.com/gallery" target="_blank">live Purtypics gallery</a> or check out our [themes](THEMES.md).

![Purtypics](.assets/purtypics.png)

## Features

- **Static Generation**: Creates fast, self-contained HTML galleries
- **Multi-Resolution**: Automatically generates multiple image sizes for optimal loading
- **Metadata Editor**: Built-in web interface for editing photo titles and descriptions
- **Video Support**: Handles videos with automatic thumbnail generation
- **EXIF Data**: Extracts and displays camera settings and location data
- **Responsive Design**: Beautiful masonry layout that works on all devices
- **[13 Built-in Themes](THEMES.md)**: From minimal to dramatic — find the right look for your gallery
- **Easy Deployment**: Deploy to any static host (rsync, S3, Cloudflare Pages)

## Quick Start

### 1. Install

<!-- begin:installer-links (auto-updated by make release) -->
| Platform | Installer |
|----------|-----------|
| **macOS (Apple Silicon)** | [purtypics_1.4.2_macOS_arm64.tar.gz](https://github.com/chrissnell/purtypics/releases/download/v1.4.2/purtypics_1.4.2_macOS_arm64.tar.gz) |
| **macOS (Intel)** | [purtypics_1.4.2_macOS_x86_64.tar.gz](https://github.com/chrissnell/purtypics/releases/download/v1.4.2/purtypics_1.4.2_macOS_x86_64.tar.gz) |
| **Linux (x86_64 deb)** | [purtypics_1.4.2_linux_amd64.deb](https://github.com/chrissnell/purtypics/releases/download/v1.4.2/purtypics_1.4.2_linux_amd64.deb) |
| **Linux (x86_64 rpm)** | [purtypics_1.4.2_linux_amd64.rpm](https://github.com/chrissnell/purtypics/releases/download/v1.4.2/purtypics_1.4.2_linux_amd64.rpm) |
| **Linux (arm64 deb)** | [purtypics_1.4.2_linux_arm64.deb](https://github.com/chrissnell/purtypics/releases/download/v1.4.2/purtypics_1.4.2_linux_arm64.deb) |
| **Windows** | [purtypics-installer.exe](https://github.com/chrissnell/purtypics/releases/download/v1.4.2/purtypics-installer.exe) |
<!-- end:installer-links -->

See all platforms on the [releases page](https://github.com/chrissnell/purtypics/releases).

### 2. Launch the editor and generate

```bash
purtypics edit ~/photos --output ./gallery
```

The `--output` flag sets where the generated site goes. The editor opens a browser UI where you can organize albums, tag favorites, set cover photos, and choose a theme — no YAML editing required. When everything looks right, click **Generate Gallery** to build the site, then **Deploy** to push it live.

## Installation

The download links above are the quickest way to get started. You can also install from source or via package managers.

### From Source

```bash
git clone https://github.com/chrissnell/purtypics.git
cd purtypics
make build
```

## Usage

### Gallery Structure

Purtypics expects your photos to be organized in directories:

```
~/photos/
├── vacation-2024/
│   ├── IMG_001.jpg
│   ├── IMG_002.jpg
│   └── VIDEO_001.mp4
├── family-reunion/
│   ├── DSC_001.jpg
│   └── DSC_002.jpg
└── gallery.yaml      # Auto-generated metadata file
```

Each subdirectory becomes an album in your gallery.

### Metadata Editor (Recommended)

The easiest way to set up and generate your gallery is the built-in editor:

```bash
purtypics edit ~/photos --output ./gallery
```

This opens a browser UI where you can:
- Add titles and descriptions to photos
- Mark favorites and hide photos
- Set album cover photos
- Choose a theme from the Gallery Settings tab
- Edit gallery title and description
- **Generate your gallery** with the Generate Gallery button

All changes are saved to `gallery.yaml`. When you're happy with everything, click **Generate Gallery** — the site is built to the `--output` directory.

### Command-Line Generation

You can also generate without the editor:

```bash
purtypics generate
purtypics generate -s ~/photos --output ./gallery
```

This scans for photos and videos, generates optimized thumbnails, and creates the static site.

### Advanced Options

#### Gallery Configuration

You can also hand-edit `gallery.yaml` to customize your gallery:

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

#### Themes

Purtypics ships with **13 built-in themes**. Set the theme in `gallery.yaml`:

```yaml
theme: darkroom
```

| Theme | Description |
|-------|-------------|
| `aperture` | Clean, modern cards with a bold accent divider |
| `atelier` | Warm, artisanal feel with creamy background and soft shadows |
| `blueprint` | Deep blue background with white grid lines and drafting-table aesthetic |
| `brutalist` | High-contrast black and white, bold type, zero decoration |
| `coast` | Breezy coastal palette with sandy tones and ocean-blue accents |
| `darkroom` | Dark background with warm amber accents — photos glow |
| `default` | Soft pastel palette with dotted borders and masonry layout |
| `ember` | Warm, fiery palette with deep reds and glowing orange highlights |
| `kyoto` | Japanese-inspired minimalism with muted earth tones |
| `mono` | Monochrome elegance — grayscale UI with clean lines |
| `nordic` | Scandinavian calm — whitespace, cool blue-gray, elegant serif headings |
| `polaroid` | Nostalgic instant-photo cards with handwriting font and slight rotation |
| `salon` | Gallery-wall presentation with rich dark background and elegant framing |

See the **[Theme Gallery](THEMES.md)** for screenshots of every theme.

You can also select a theme from the **Gallery Settings** tab in the editor (`purtypics edit`).

#### Custom Themes

Create a custom theme by adding a directory under `themes/` in your gallery source:

```
~/photos/
├── vacation-2024/
├── gallery.yaml
└── themes/
    └── mytheme/
        ├── css/
        │   └── gallery.css
        └── templates/       # optional
            ├── base.html
            ├── index.html
            └── album.html
```

Then set `theme: mytheme` in `gallery.yaml`.

A theme only needs to include files you want to override — anything missing falls back to the built-in default. For example, a CSS-only theme just needs `css/gallery.css`. Look at the [default theme](pkg/gallery/assets/themes/default/) as a reference for available template variables and CSS classes.

**Theme search order:**

1. **System-installed** — `/usr/share/purtypics/themes/` (Linux), `/usr/local/share/purtypics/themes/` or `/opt/homebrew/share/purtypics/themes/` (macOS), `%ProgramData%\Purtypics\themes\` (Windows). Installed by `make install` or platform packages.
2. **Local** — `themes/` directory alongside your photos
3. **Built-in** — themes compiled into the binary

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

#### S3 Deployment

Deploy to an S3 bucket with incremental sync — only changed files are uploaded and stale files are deleted. Deduplication uses MD5 ETag comparison.

```yaml
s3:
  bucket: my-photo-gallery
  region: us-east-1
  # Optional settings:
  cloudfront_id: E27EXAMPLE51Z    # invalidate CloudFront cache after deploy
  cache_control: "max-age=31536000"
  storage_class: STANDARD          # STANDARD, STANDARD_IA, GLACIER, etc.
  acl: public-read                 # public-read, private, etc.
```

AWS credentials are resolved via the standard SDK chain: environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`), `~/.aws/credentials`, or IAM role.

```bash
purtypics deploy --target s3
purtypics deploy --target s3 --dry-run  # test connection only
```

#### Cloudflare Pages

Deploy to Cloudflare Pages with hash-based deduplication. Large files (>25 MiB) are automatically uploaded to R2 if enabled.

```yaml
cloudflare:
  project: my-gallery
  account_id: your-account-id
  auto_create: true               # create project if it doesn't exist
  branch: main                    # optional branch name
  r2:                             # optional: handle files >25 MiB
    enabled: true
    bucket: my-gallery-assets     # defaults to {project}-assets
    custom_domain: assets.example.com
```

Requires `CLOUDFLARE_API_TOKEN` environment variable. For R2, also set `R2_ACCESS_KEY_ID` and `R2_SECRET_ACCESS_KEY`.

```bash
purtypics deploy --target cloudflare
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
