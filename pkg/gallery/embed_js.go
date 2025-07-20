package gallery

import (
	_ "embed"
)

//go:embed js/bundle.js
var galleryJSContent string

// GetGalleryJS returns the bundled JavaScript
func GetGalleryJS() string {
	return galleryJSContent
}