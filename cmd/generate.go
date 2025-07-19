package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjs/purtypics/pkg/gallery"
	"github.com/spf13/cobra"
)

var (
	generateSource   string
	generateOutput   string
	generateBaseURL  string
	generateTitle    string
	generateMetadata string
	generateVerbose  bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a static photo gallery",
	Long: `Generate a static photo gallery from a directory of photos.
	
The generated gallery includes thumbnails, optimized images, and an HTML viewer.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate required flags
		if generateSource == "" {
			return fmt.Errorf("source directory is required")
		}
		if generateOutput == "" {
			return fmt.Errorf("output directory is required")
		}

		// Ensure source directory exists
		if _, err := os.Stat(generateSource); os.IsNotExist(err) {
			return fmt.Errorf("source directory does not exist: %s", generateSource)
		}

		// Create gallery generator
		generator := gallery.NewGenerator(generateSource, generateOutput, generateTitle, generateBaseURL, generateVerbose)
		
		// Set metadata path if provided
		if generateMetadata != "" {
			metadataPath := generateMetadata
			if !filepath.IsAbs(metadataPath) {
				metadataPath = filepath.Join(generateSource, metadataPath)
			}
			generator.MetadataPath = metadataPath
		}

		// Generate the gallery
		if err := generator.Generate(); err != nil {
			return fmt.Errorf("failed to generate gallery: %w", err)
		}

		fmt.Printf("Gallery generated successfully at %s\n", generateOutput)
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&generateSource, "source", "s", "", "Source directory containing photos (required)")
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "Output directory for the gallery (required)")
	generateCmd.Flags().StringVar(&generateBaseURL, "baseurl", "", "Base URL for the gallery (e.g., https://example.com/photos)")
	generateCmd.Flags().StringVar(&generateTitle, "title", "Photo Gallery", "Title for the gallery")
	generateCmd.Flags().StringVar(&generateMetadata, "metadata", "gallery.yaml", "Path to metadata file (relative to source or absolute)")
	generateCmd.Flags().BoolVarP(&generateVerbose, "verbose", "v", false, "Enable verbose output")

	generateCmd.MarkFlagRequired("source")
	generateCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(generateCmd)
}