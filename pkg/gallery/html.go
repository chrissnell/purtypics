package gallery

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Breadcrumb represents a single breadcrumb item
type Breadcrumb struct {
	Title string
	URL   string
}

// HTMLGenerator creates static HTML galleries
type HTMLGenerator struct {
	OutputPath    string
	SiteTitle     string
	BaseURL       string
	ShowLocations bool
}

// Generate creates the complete static site
func (g *HTMLGenerator) Generate(albums []Album) error {
	// Create directory structure
	if err := g.createDirectories(); err != nil {
		return err
	}

	// Copy static assets
	if err := g.writeStaticAssets(); err != nil {
		return err
	}

	// Generate album pages
	for _, album := range albums {
		if err := g.generateAlbumPage(album); err != nil {
			return err
		}
	}

	// Generate index page
	return g.generateIndexPage(albums)
}

func (g *HTMLGenerator) createDirectories() error {
	dirs := []string{
		"albums",
		"css",
		"js",
		"static",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(g.OutputPath, dir), 0755); err != nil {
			return err
		}
	}
	return nil
}

func (g *HTMLGenerator) writeStaticAssets() error {
	// Write CSS
	cssPath := filepath.Join(g.OutputPath, "css", "style.css")
	if err := os.WriteFile(cssPath, []byte(galleryCSS), 0644); err != nil {
		return err
	}

	// Write JavaScript
	jsPath := filepath.Join(g.OutputPath, "js", "gallery.js")
	if err := os.WriteFile(jsPath, []byte(galleryJS), 0644); err != nil {
		return err
	}

	return nil
}

func (g *HTMLGenerator) generateIndexPage(albums []Album) error {
	tmpl, err := template.New("base").Parse(baseTemplate)
	if err != nil {
		return err
	}

	contentTmpl, err := template.New("content").Parse(indexContent)
	if err != nil {
		return err
	}

	var contentBuf strings.Builder
	contentData := struct {
		Albums  []Album
		BaseURL string
	}{
		Albums:  albums,
		BaseURL: ".",
	}

	if err := contentTmpl.Execute(&contentBuf, contentData); err != nil {
		return err
	}

	pageData := struct {
		Title       string
		SiteTitle   string
		BaseURL     string
		Content     template.HTML
		Breadcrumbs []Breadcrumb
	}{
		Title:       g.SiteTitle, // Just the gallery title for index page
		SiteTitle:   g.SiteTitle,
		BaseURL:     ".",
		Content:     template.HTML(contentBuf.String()),
		Breadcrumbs: nil, // No breadcrumbs on index page
	}

	indexPath := filepath.Join(g.OutputPath, "index.html")
	file, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, pageData)
}

func (g *HTMLGenerator) generateAlbumPage(album Album) error {
	tmpl, err := template.New("base").Parse(baseTemplate)
	if err != nil {
		return err
	}

	contentTmpl, err := template.New("content").Parse(albumContent)
	if err != nil {
		return err
	}

	// Add aspect ratio to photos
	for i := range album.Photos {
		photo := &album.Photos[i]
		if photo.Width > 0 && photo.Height > 0 {
			ratio := float64(photo.Width) / float64(photo.Height)
			if ratio > 1.5 {
				photo.AspectRatio = "wide"
			} else if ratio < 0.8 {
				photo.AspectRatio = "tall"
			} else {
				photo.AspectRatio = "normal"
			}
		}
	}

	// Check if any photos have GPS data AND if showing locations is enabled
	hasGPS := false
	if g.ShowLocations {
		for _, photo := range album.Photos {
			if photo.EXIF != nil && photo.EXIF.GPS != nil {
				hasGPS = true
				break
			}
		}
	}

	var contentBuf strings.Builder
	contentData := struct {
		Album   Album
		BaseURL string
		HasGPS  bool
	}{
		Album:   album,
		BaseURL: "..",
		HasGPS:  hasGPS,
	}

	if err := contentTmpl.Execute(&contentBuf, contentData); err != nil {
		return err
	}

	pageData := struct {
		Title       string
		SiteTitle   string
		BaseURL     string
		Content     template.HTML
		Breadcrumbs []Breadcrumb
	}{
		Title:       album.Title,
		SiteTitle:   g.SiteTitle,
		BaseURL:     "..",
		Content:     template.HTML(contentBuf.String()),
		Breadcrumbs: []Breadcrumb{
			{Title: g.SiteTitle, URL: "../"},
			{Title: album.Title, URL: ""}, // Current page, no link
		},
	}

	albumPath := filepath.Join(g.OutputPath, "albums", fmt.Sprintf("%s.html", album.ID))
	file, err := os.Create(albumPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, pageData)
}
