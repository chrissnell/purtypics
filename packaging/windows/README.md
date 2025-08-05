# Windows Packaging for Purtypics

This directory contains scripts and configuration for building Windows binaries and installers for Purtypics.

## Cross-Platform Installer Generation

You can create Windows installers on macOS or Linux using several tools:

### 1. NSIS (Recommended for cross-platform)
- **Install on macOS**: `brew install nsis`
- **Install on Linux**: `apt-get install nsis` (Debian/Ubuntu) or `yum install nsis` (Fedora)
- **Build**: `make nsis`

### 2. WiX with msitools (Creates MSI files)
- **Install on macOS**: `brew install msitools`
- **Install on Linux**: `apt-get install msitools` (Debian/Ubuntu)
- **Build**: `make msi`

### 3. Self-Extracting ZIP (Simplest)
- No additional tools required
- **Build**: `make sfx`

## Prerequisites

1. **Go**: Install Go 1.24 or later from https://go.dev/dl/
2. **Installer tools** (choose one):
   - NSIS: Cross-platform, creates .exe installers
   - msitools: Cross-platform, creates .msi installers
   - Inno Setup: Windows-only, creates .exe installers

## Building

### From macOS/Linux:

```bash
cd packaging/windows

# Build Windows binary only
make build

# Create self-extracting ZIP (no tools required)
make sfx

# Create NSIS installer (requires NSIS)
make nsis

# Create MSI installer (requires msitools)
make msi

# Build binary and create all packages
make all
```

### From project root:

```bash
# Build Windows binary
make build-all

# Create distribution packages
make dist
```

## Files

- `Makefile` - Cross-platform build configuration
- `purtypics.iss` - Inno Setup script (Windows-only)
- `purtypics.nsi` - NSIS script (cross-platform)
- `purtypics.wxs` - WiX configuration for MSI (cross-platform with msitools)
- `build.bat` - Legacy Windows batch script
- `README.md` - This file

## Installer Features

All installers provide:
- Installation to Program Files
- Optional PATH environment variable update
- Start menu shortcuts
- Desktop shortcut
- Uninstaller
- Version information in Windows Programs list

## FFmpeg Installation for Windows

Purtypics requires FFmpeg for video thumbnail extraction. Windows users have several options:

### Option 1: Download Pre-built FFmpeg (Recommended)
1. Visit https://www.gyan.dev/ffmpeg/builds/
2. Download the "full" build (not essentials)
3. Extract the ZIP file to `C:fmpeg`
4. Add `C:fmpegin` to your system PATH:
   - Right-click "This PC" → Properties → Advanced system settings
   - Click "Environment Variables"
   - Under System variables, find "Path" and click Edit
   - Click New and add `C:fmpegin`
   - Click OK to save

### Option 2: Using Chocolatey
If you have Chocolatey installed:
```powershell
choco install ffmpeg
```

### Option 3: Using Scoop
If you have Scoop installed:
```powershell
scoop install ffmpeg
```

### Option 4: Using winget
```
winget install ffmpeg
```

### Verify Installation
Open a new Command Prompt or PowerShell and run:
```
ffmpeg -version
```

## Notes

- All installers target Windows 64-bit systems
- FFmpeg must be installed separately for video support (see above)
- PATH modifications may require administrator privileges
- The MSI installer created with msitools is more basic than a full WiX toolset build
- If ffmpeg is not found, video thumbnails will not be generated but the application will continue to work
