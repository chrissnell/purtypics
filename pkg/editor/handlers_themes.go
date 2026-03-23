package editor

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/gallery"
)

// handleThemes returns available theme names (built-in + system-installed + local)
func (s *Server) handleThemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seen := make(map[string]bool)
	var themes []string

	// Built-in embedded themes
	embeddedFS := gallery.EmbeddedThemesFS()
	if entries, err := fs.ReadDir(embeddedFS, "."); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				themes = append(themes, e.Name())
				seen[e.Name()] = true
			}
		}
	}

	// System-installed themes (platform-specific directories)
	for _, dir := range gallery.SystemThemesDirs() {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, e := range entries {
				if e.IsDir() && !seen[e.Name()] {
					themes = append(themes, e.Name())
					seen[e.Name()] = true
				}
			}
		}
	}

	// Local themes from <source>/themes/
	localThemesDir := filepath.Join(s.SourcePath, "themes")
	if entries, err := os.ReadDir(localThemesDir); err == nil {
		for _, e := range entries {
			if e.IsDir() && !seen[e.Name()] {
				themes = append(themes, e.Name())
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(themes)
}
