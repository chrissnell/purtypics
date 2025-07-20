package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
	Use:   "edit",
	Short: "Launch the web-based metadata editor",
	Long: `Launch a web-based editor for editing photo metadata.
	
The editor allows you to:
  - View and edit photo titles and descriptions
  - Mark photos as favorites
  - Hide photos from the gallery
  - Save metadata to a YAML file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate required flags
		if editSource == "" {
			return fmt.Errorf("source directory is required")
		}
		if editOutput == "" {
			editOutput = editSource // Default output to source if not specified
		}

		// Ensure source directory exists
		if _, err := os.Stat(editSource); os.IsNotExist(err) {
			return fmt.Errorf("source directory does not exist: %s", editSource)
		}

		// Determine metadata file path
		metadataPath := editMetadata
		if !filepath.IsAbs(metadataPath) {
			metadataPath = filepath.Join(editSource, metadataPath)
		}

		// Create and start the editor server
		server := editor.NewServer(editSource, metadataPath, editPort)
		server.OutputPath = editOutput
		
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
			openBrowser(url)
		}

		fmt.Println("Press Ctrl+C to stop the server")
		
		return server.StartWithListener(listener)
	},
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		fmt.Printf("Please open %s in your browser\n", url)
		return
	}
	
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
		fmt.Printf("Please open %s in your browser\n", url)
	}
}

func init() {
	editCmd.Flags().StringVarP(&editSource, "source", "s", "", "Source directory containing photos (required)")
	editCmd.Flags().StringVarP(&editOutput, "output", "o", "", "Output directory for metadata file (defaults to source)")
	editCmd.Flags().StringVar(&editMetadata, "metadata", "gallery.yaml", "Path to metadata file (relative to output or absolute)")
	editCmd.Flags().IntVarP(&editPort, "port", "p", 8080, "Port to run the editor server on")
	editCmd.Flags().BoolVar(&editNoBrowser, "no-browser", false, "Don't open browser automatically")

	editCmd.MarkFlagRequired("source")

	rootCmd.AddCommand(editCmd)
}