package gallery

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed all:assets
var assetsFS embed.FS

const defaultTheme = "default"

// EmbeddedThemesFS returns the embedded themes directory for listing built-in themes.
func EmbeddedThemesFS() fs.FS {
	sub, _ := fs.Sub(assetsFS, "assets/themes")
	return sub
}

// SystemThemesDirs returns the platform-specific system-wide themes directories,
// ordered by priority (first match wins).
//
//	macOS:   /usr/local/share/purtypics/themes, /usr/share/purtypics/themes
//	Linux:   /usr/share/purtypics/themes
//	Windows: %ProgramData%\Purtypics\themes
func SystemThemesDirs() []string {
	switch runtime.GOOS {
	case "windows":
		if pd := os.Getenv("ProgramData"); pd != "" {
			return []string{filepath.Join(pd, "Purtypics", "themes")}
		}
		return []string{`C:\ProgramData\Purtypics\themes`}
	case "darwin":
		return []string{
			"/opt/homebrew/share/purtypics/themes",
			"/usr/local/share/purtypics/themes",
		}
	default:
		return []string{"/usr/share/purtypics/themes"}
	}
}

// ThemeFS provides access to theme files with fallback to the embedded default theme.
//
// Resolution order for a given theme name:
//  1. System-installed themes: /usr/share/purtypics/themes/<name>/ (or platform equivalent)
//  2. Local themes: <sourcePath>/themes/<name>/
//  3. Embedded themes compiled into the binary
type ThemeFS struct {
	themeName string
	userFS    fs.FS // highest-priority on-disk theme (may be nil)
	defaultFS fs.FS // fallback theme (embedded default or embedded named theme)
}

// NewThemeFS creates a ThemeFS for the given theme name.
// sourcePath is the gallery source directory where local themes live under themes/.
func NewThemeFS(themeName, sourcePath string) (*ThemeFS, error) {
	if themeName == "" {
		themeName = defaultTheme
	}

	// Embedded default is always the ultimate fallback
	embeddedDefault, err := fs.Sub(assetsFS, "assets/themes/default")
	if err != nil {
		return nil, fmt.Errorf("failed to access embedded default theme: %w", err)
	}

	t := &ThemeFS{
		themeName: themeName,
		defaultFS: embeddedDefault,
	}

	if themeName == defaultTheme {
		return t, nil
	}

	// For non-default themes, resolve in priority order:
	//   1. System-installed (/usr/share/purtypics/themes/<name>/)
	//   2. Local (<sourcePath>/themes/<name>/)
	//   3. Embedded (assets/themes/<name>/)

	// Check system-installed themes
	for _, dir := range SystemThemesDirs() {
		sysDir := filepath.Join(dir, themeName)
		if info, err := os.Stat(sysDir); err == nil && info.IsDir() {
			t.userFS = os.DirFS(sysDir)
			return t, nil
		}
	}

	// Check local themes in source directory
	if sourcePath != "" {
		localDir := filepath.Join(sourcePath, "themes", themeName)
		if info, err := os.Stat(localDir); err == nil && info.IsDir() {
			t.userFS = os.DirFS(localDir)
			return t, nil
		}
	}

	// Check embedded themes
	if embeddedTheme, err := fs.Sub(assetsFS, "assets/themes/"+themeName); err == nil {
		if _, err := fs.ReadDir(embeddedTheme, "."); err == nil {
			t.defaultFS = embeddedTheme
			return t, nil
		}
	}

	return nil, fmt.Errorf("theme %q not found (checked system, local, and built-in themes)", themeName)
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

	if rdr, ok := o.lower.(fs.ReadDirFS); ok {
		if dirEntries, err := rdr.ReadDir(name); err == nil {
			for _, e := range dirEntries {
				entries[e.Name()] = e
			}
		}
	}

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
