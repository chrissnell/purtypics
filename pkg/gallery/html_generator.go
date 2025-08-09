package gallery

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// GalleryData represents gallery data for templates
type GalleryData struct {
	Title       string
	Description string
	Author      string
	Copyright   string
	Albums      []Album
}

// HTMLTemplateData represents the data passed to HTML templates
type HTMLTemplateData struct {
	Title       string
	Description string
	BasePath    string
	Content     template.HTML
	Gallery     *GalleryData
	Album       *Album
	Version     string
	CommitHash  string
}

// GenerateHTMLFromTemplates generates the gallery HTML using the embedded templates
func (g *Generator) GenerateHTMLFromTemplates(albums []Album) error {
	// Get the template filesystem
	templateFS, err := GetTemplateFS()
	if err != nil {
		return fmt.Errorf("failed to get template filesystem: %w", err)
	}

	// Parse all templates
	tmpl, err := template.ParseFS(templateFS, "*.html")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Copy static assets (CSS, JS)
	if err := g.copyStaticAssets(); err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	// Create gallery data
	galleryData := &GalleryData{
		Title:       g.SiteTitle,
		Description: "",
		Albums:      albums,
	}

	// Use metadata if available
	if g.metadata != nil {
		if g.metadata.Title != "" {
			galleryData.Title = g.metadata.Title
		}
		galleryData.Description = g.metadata.Description
		galleryData.Author = g.metadata.Author
		galleryData.Copyright = g.metadata.Copyright
	}

	// Generate index page
	if err := g.generateIndexPage(tmpl, galleryData); err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Generate album pages
	for i := range albums {
		if err := g.generateAlbumPage(tmpl, &albums[i], galleryData); err != nil {
			return fmt.Errorf("failed to generate album page for %s: %w", albums[i].ID, err)
		}
	}

	return nil
}

// copyStaticAssets copies CSS and JS files to the output directory
func (g *Generator) copyStaticAssets() error {
	staticFS, err := GetStaticFS()
	if err != nil {
		return err
	}

	// Create css and js directories
	cssDir := filepath.Join(g.OutputPath, "css")
	jsDir := filepath.Join(g.OutputPath, "js")

	if err := os.MkdirAll(cssDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(jsDir, 0755); err != nil {
		return err
	}

	// Copy CSS files
	if err := fs.WalkDir(staticFS, "css", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		src, err := staticFS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		destPath := filepath.Join(g.OutputPath, path)
		dest, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		return err
	}); err != nil {
		return fmt.Errorf("failed to copy CSS files: %w", err)
	}

	// Copy JS files
	if err := fs.WalkDir(staticFS, "js", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		src, err := staticFS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		destPath := filepath.Join(g.OutputPath, path)
		dest, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		return err
	}); err != nil {
		return fmt.Errorf("failed to copy JS files: %w", err)
	}

	return nil
}

// generateIndexPage generates the gallery index page
func (g *Generator) generateIndexPage(tmpl *template.Template, galleryData *GalleryData) error {

	// Render the index content
	var contentBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&contentBuf, "index.html", HTMLTemplateData{
		Gallery: galleryData,
	}); err != nil {
		return fmt.Errorf("failed to render index content: %w", err)
	}

	// Render the full page
	var pageBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&pageBuf, "base.html", HTMLTemplateData{
		Title:       galleryData.Title,
		Description: galleryData.Description,
		BasePath:    ".",
		Content:     template.HTML(contentBuf.String()),
		Version:     g.Version,
		CommitHash:  g.CommitHash,
	}); err != nil {
		return fmt.Errorf("failed to render index page: %w", err)
	}

	// Write to file
	indexPath := filepath.Join(g.OutputPath, "index.html")
	return os.WriteFile(indexPath, pageBuf.Bytes(), 0644)
}

// generateAlbumPage generates a single album page
func (g *Generator) generateAlbumPage(tmpl *template.Template, album *Album, galleryData *GalleryData) error {
	// Create album directory
	albumDir := filepath.Join(g.OutputPath, album.ID)
	if err := os.MkdirAll(albumDir, 0755); err != nil {
		return err
	}

	// Render the album content
	var contentBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&contentBuf, "album.html", HTMLTemplateData{
		Album:   album,
		Gallery: galleryData,
	}); err != nil {
		return fmt.Errorf("failed to render album content: %w", err)
	}

	// Render the full page
	var pageBuf bytes.Buffer
	title := fmt.Sprintf("%s - %s", album.Title, g.SiteTitle)
	// Use metadata title if available
	if g.metadata != nil && g.metadata.Title != "" {
		title = fmt.Sprintf("%s - %s", album.Title, g.metadata.Title)
	}
	if err := tmpl.ExecuteTemplate(&pageBuf, "base.html", HTMLTemplateData{
		Title:       title,
		Description: string(album.Description),
		BasePath:    "..",
		Content:     template.HTML(contentBuf.String()),
		Version:     g.Version,
		CommitHash:  g.CommitHash,
	}); err != nil {
		return fmt.Errorf("failed to render album page: %w", err)
	}

	// Write to file
	albumIndexPath := filepath.Join(albumDir, "index.html")
	return os.WriteFile(albumIndexPath, pageBuf.Bytes(), 0644)
}