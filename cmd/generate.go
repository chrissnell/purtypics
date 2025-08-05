package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/common"
	"github.com/cjs/purtypics/pkg/gallery"
	"github.com/cjs/purtypics/pkg/metadata"
	"github.com/spf13/cobra"
)

var (
	generateSource   string
	generateOutput   string
	generateMetadata string
	generateTitle    string
	generateVerbose  bool
)

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate a static photo gallery",
	Long: `Generate a static photo gallery from a directory of photos.
	
The generated gallery includes thumbnails, optimized images, and an HTML viewer.

Usage:
  purtypics generate                    # Generate in current directory
  purtypics generate /path/to/gallery   # Generate in specified directory
  purtypics generate -s /photos -o /web # Use explicit source and output paths

When using explicit flags (-s and -o), the command works with any directory structure.
Otherwise, it expects a gallery directory with photos, gallery.yaml, and output/ subdirectory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var sourcePath, outputPath, metadataPath string
		
		// Check if using explicit flags
		if generateSource != "" || generateOutput != "" {
			// Legacy mode with explicit paths
			if generateSource == "" {
				return fmt.Errorf("source directory is required when using explicit paths")
			}
			if generateOutput == "" {
				return fmt.Errorf("output directory is required when using explicit paths")
			}
			
			sourcePath = generateSource
			outputPath = generateOutput
			
			// Handle metadata path
			if generateMetadata != "" {
				metadataPath = common.ResolvePath(generateMetadata, sourcePath)
			} else {
				metadataPath = filepath.Join(sourcePath, "gallery.yaml")
			}
		} else {
			// New mode matching web UI behavior
			galleryPath := "."
			if len(args) > 0 {
				galleryPath = args[0]
			}
			
			sourcePath = galleryPath
			// Default output to sibling gallery directory
			if filepath.IsAbs(galleryPath) {
				// For absolute paths like /photos -> /gallery
				outputPath = filepath.Join(filepath.Dir(galleryPath), "gallery")
			} else {
				// For relative paths like . -> gallery
				outputPath = "gallery"
			}
			metadataPath = filepath.Join(galleryPath, "gallery.yaml")
		}

		// Ensure source directory exists
		if err := common.ValidateDirectory(sourcePath); err != nil {
			return err
		}

		// Load metadata to get title
		var title string
		if generateTitle != "" {
			title = generateTitle
		} else {
			if data, err := metadata.Load(metadataPath); err == nil && data != nil {
				title = data.Title
			}
			if title == "" {
				title = "Photo Gallery"
			}
		}

		// Create a progress callback that prints to console
		progressCallback := func(current, total int, message string) {
			if total > 0 {
				progress := (current * 100) / total
				// Clear the line and print progress
				fmt.Printf("\r\033[K[%3d%%] %s", progress, message)
				if current == total {
					fmt.Println() // New line when complete
				}
			} else {
				// If no total, just print the message
				fmt.Printf("\r\033[K%s", message)
			}
		}

		// Create gallery generator
		generator := gallery.NewGenerator(sourcePath, outputPath, title, "", generateVerbose)
		generator.MetadataPath = metadataPath
		generator.ProgressCallback = progressCallback

		fmt.Printf("Generating gallery from %s...\n", sourcePath)
		
		// Generate the gallery
		if err := generator.Generate(); err != nil {
			return fmt.Errorf("failed to generate gallery: %w", err)
		}

		fmt.Printf("\nGallery generated successfully at %s\n", outputPath)
		
		// Only suggest edit command if using the new mode
		if generateSource == "" && generateOutput == "" {
			fmt.Printf("To view the gallery, run: purtypics edit %s\n", sourcePath)
		}
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&generateSource, "source", "s", "", "Source directory containing photos (overrides default behavior)")
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "Output directory for the gallery (overrides default behavior)")
	generateCmd.Flags().StringVar(&generateMetadata, "metadata", "", "Path to metadata file (default: gallery.yaml in source)")
	generateCmd.Flags().StringVar(&generateTitle, "title", "", "Title for the gallery (overrides metadata)")
	generateCmd.Flags().BoolVarP(&generateVerbose, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(generateCmd)
}