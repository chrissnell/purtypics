package editor

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/gallery"
)

// handleThemes returns available theme names (built-in + user-provided)
func (s *Server) handleThemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seen := make(map[string]bool)
	var themes []string

	// Built-in themes from embedded assets
	embeddedFS := gallery.EmbeddedThemesFS()
	if entries, err := fs.ReadDir(embeddedFS, "."); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				themes = append(themes, e.Name())
				seen[e.Name()] = true
			}
		}
	}

	// User-provided themes from <source>/themes/
	userThemesDir := filepath.Join(s.SourcePath, "themes")
	if entries, err := os.ReadDir(userThemesDir); err == nil {
		for _, e := range entries {
			if e.IsDir() && !seen[e.Name()] {
				themes = append(themes, e.Name())
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(themes)
}
