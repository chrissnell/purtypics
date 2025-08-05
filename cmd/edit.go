package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/common"
	"github.com/cjs/purtypics/pkg/editor"
	"github.com/spf13/cobra"
)

var (
	editSource    string
	editOutput    string
	editMetadata  string
	editPort      int
	editNoBrowser bool
)

var editCmd = &cobra.Command{
	Use:   "edit [path]",
	Short: "Launch the web-based metadata editor",
	Long: `Launch a web-based editor for editing photo metadata.
	
The editor allows you to:
  - View and edit photo titles and descriptions
  - Mark photos as favorites
  - Hide photos from the gallery
  - Save metadata to a YAML file
  - Generate the gallery HTML and thumbnails

Usage:
  purtypics edit                    # Edit gallery in current directory
  purtypics edit /path/to/gallery   # Edit gallery in specified directory
  purtypics edit -s /photos -o /web # Use explicit source and output paths

The metadata file (gallery.yaml) is saved in the source directory by default.
When you click "Generate" in the editor, the gallery HTML and thumbnails are
created in the output directory (-o flag) or in a default location.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var sourcePath, outputPath, metadataPath string

		// Check if using explicit flags
		if editSource != "" || editOutput != "" {
			// Legacy mode with explicit paths
			if editSource == "" {
				return fmt.Errorf("source directory is required when using explicit paths")
			}
			
			sourcePath = editSource
			outputPath = editOutput
			if outputPath == "" {
				// Default output to sibling gallery directory
				if filepath.IsAbs(sourcePath) {
					// For absolute paths like /photos -> /gallery
					outputPath = filepath.Join(filepath.Dir(sourcePath), "gallery")
				} else {
					// For relative paths
					outputPath = "gallery"
				}
			}
			metadataPath = common.ResolvePath(editMetadata, sourcePath)
		} else {
			// New mode matching generate behavior
			galleryPath := "."
			if len(args) > 0 {
				galleryPath = args[0]
			}
			
			sourcePath = galleryPath
			// Default output to ./gallery
			outputPath = "gallery"
			metadataPath = filepath.Join(galleryPath, "gallery.yaml")
		}

		// Ensure source directory exists
		if err := common.ValidateDirectory(sourcePath); err != nil {
			return err
		}

		// Create and start the editor server
		server := editor.NewServer(sourcePath, metadataPath, editPort)
		server.OutputPath = outputPath
		
		// Get the actual port the server will listen on
		actualPort, listener, err := server.GetActualPort()
		if err != nil {
			return err
		}
		
		fmt.Printf("Starting metadata editor on http://localhost:%d\n", actualPort)
		
		// Open browser if not disabled
		if !editNoBrowser {
			url := fmt.Sprintf("http://localhost:%d", actualPort)
			fmt.Printf("Opening browser at %s\n", url)
			if err := common.OpenBrowser(url); err != nil {
				fmt.Printf("Failed to open browser: %v\n", err)
				fmt.Printf("Please open %s in your browser\n", url)
			}
		}

		fmt.Println("Press Ctrl+C to stop the server")
		
		return server.StartWithListener(listener)
	},
}

func init() {
	editCmd.Flags().StringVarP(&editSource, "source", "s", "", "Source directory containing photos (overrides default behavior)")
	editCmd.Flags().StringVarP(&editOutput, "output", "o", "", "Output directory for generated gallery HTML and thumbnails")
	editCmd.Flags().StringVar(&editMetadata, "metadata", "gallery.yaml", "Path to metadata file (relative to source or absolute)")
	editCmd.Flags().IntVarP(&editPort, "port", "p", 0, "Port to run the editor server on (0 for auto-assign)")
	editCmd.Flags().BoolVar(&editNoBrowser, "no-browser", false, "Don't open browser automatically")

	rootCmd.AddCommand(editCmd)
}