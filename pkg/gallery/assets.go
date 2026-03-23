package gallery

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:assets
var assetsFS embed.FS

const defaultTheme = "default"

// ThemeFS provides access to theme files with fallback to the embedded default theme.
// Resolution order for a given theme name:
//   1. User theme directory: <sourcePath>/themes/<name>/
//   2. Embedded default theme (compiled into binary)
type ThemeFS struct {
	themeName string
	userFS    fs.FS // user-provided theme on disk (may be nil)
	defaultFS fs.FS // embedded default theme (always present)
}

// NewThemeFS creates a ThemeFS for the given theme name.
// sourcePath is the gallery source directory where user themes live under themes/.
func NewThemeFS(themeName, sourcePath string) (*ThemeFS, error) {
	if themeName == "" {
		themeName = defaultTheme
	}

	// Embedded default is always available
	defaultFS, err := fs.Sub(assetsFS, "assets/themes/default")
	if err != nil {
		return nil, fmt.Errorf("failed to access embedded default theme: %w", err)
	}

	t := &ThemeFS{
		themeName: themeName,
		defaultFS: defaultFS,
	}

	// Check for user-provided theme on disk
	if sourcePath != "" {
		userThemeDir := filepath.Join(sourcePath, "themes", themeName)
		if info, err := os.Stat(userThemeDir); err == nil && info.IsDir() {
			t.userFS = os.DirFS(userThemeDir)
		} else if themeName != defaultTheme {
			return nil, fmt.Errorf("theme %q not found at %s", themeName, userThemeDir)
		}
	}

	return t, nil
}

// GetTemplateFS returns an fs.FS rooted at the templates/ subdirectory,
// overlaying user theme templates on top of the default.
func (t *ThemeFS) GetTemplateFS() (fs.FS, error) {
	return t.sub("templates")
}

// GetStaticFS returns an fs.FS containing css/ and js/ subdirectories,
// overlaying user theme statics on top of the default.
func (t *ThemeFS) GetStaticFS() (fs.FS, error) {
	return &overlayFS{upper: t.userFS, lower: t.defaultFS}, nil
}

// sub returns an overlayFS for a subdirectory path.
func (t *ThemeFS) sub(dir string) (fs.FS, error) {
	defaultSub, err := fs.Sub(t.defaultFS, dir)
	if err != nil {
		return nil, fmt.Errorf("embedded default theme missing %s/: %w", dir, err)
	}

	var userSub fs.FS
	if t.userFS != nil {
		if s, err := fs.Sub(t.userFS, dir); err == nil {
			userSub = s
		}
		// If user theme doesn't have this subdir, that's fine — fall back entirely
	}

	return &overlayFS{upper: userSub, lower: defaultSub}, nil
}

// overlayFS layers an upper filesystem over a lower one.
// Files in upper take precedence; missing files fall through to lower.
type overlayFS struct {
	upper fs.FS // may be nil
	lower fs.FS
}

func (o *overlayFS) Open(name string) (fs.File, error) {
	if o.upper != nil {
		if f, err := o.upper.Open(name); err == nil {
			return f, nil
		}
	}
	return o.lower.Open(name)
}

// ReadDir merges directory listings from both layers.
// Upper entries take precedence for duplicate names.
func (o *overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries := make(map[string]fs.DirEntry)

	// Lower first
	if rdr, ok := o.lower.(fs.ReadDirFS); ok {
		if dirEntries, err := rdr.ReadDir(name); err == nil {
			for _, e := range dirEntries {
				entries[e.Name()] = e
			}
		}
	}

	// Upper overwrites
	if o.upper != nil {
		if rdr, ok := o.upper.(fs.ReadDirFS); ok {
			if dirEntries, err := rdr.ReadDir(name); err == nil {
				for _, e := range dirEntries {
					entries[e.Name()] = e
				}
			}
		}
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("directory %q not found in theme", name)
	}

	result := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		result = append(result, e)
	}
	return result, nil
}

// ReadFile implements fs.ReadFileFS for efficient reads.
func (o *overlayFS) ReadFile(name string) ([]byte, error) {
	if o.upper != nil {
		if rdr, ok := o.upper.(fs.ReadFileFS); ok {
			if data, err := rdr.ReadFile(name); err == nil {
				return data, nil
			}
		}
		// Try Open as fallback for upper
		if f, err := o.upper.Open(name); err == nil {
			defer f.Close()
			info, err := f.Stat()
			if err == nil {
				data := make([]byte, info.Size())
				if _, err := f.Read(data); err == nil {
					return data, nil
				}
			}
		}
	}
	if rdr, ok := o.lower.(fs.ReadFileFS); ok {
		return rdr.ReadFile(name)
	}
	f, err := o.lower.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, info.Size())
	_, err = f.Read(data)
	return data, err
}
