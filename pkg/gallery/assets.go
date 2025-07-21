package gallery

import (
	"embed"
	"io/fs"
)

// Embed all assets (templates, CSS, JS) using io/fs
//
//go:embed all:assets
var assetsFS embed.FS

// Assets provides access to the embedded assets filesystem
var Assets fs.FS

func init() {
	// Strip the "assets" prefix from the embedded filesystem
	var err error
	Assets, err = fs.Sub(assetsFS, "assets")
	if err != nil {
		panic("failed to create assets sub-filesystem: " + err.Error())
	}
}

// GetTemplateFS returns the templates subdirectory
func GetTemplateFS() (fs.FS, error) {
	return fs.Sub(Assets, "templates")
}

// GetStaticFS returns the static files (css, js) subdirectory
func GetStaticFS() (fs.FS, error) {
	return Assets, nil
}