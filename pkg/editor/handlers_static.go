package editor

import (
	"net/http"
	"strings"
)

// handleIndex serves the main editor HTML
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(editorHTML))
}

// handleStatic serves static assets (CSS, JS)
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	
	switch path {
	case "editor.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(editorCSS))
	case "editor.js":
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(editorJS))
	default:
		http.NotFound(w, r)
	}
}