package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cjs/purtypics/pkg/editor"
	"github.com/cjs/purtypics/pkg/gallery"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	
	// Remove command from args for flag parsing
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	
	// Parse flags after command is removed
	flag.Parse()

	switch command {
	case "generate":
		runGenerate()
	case "edit":
		runEdit()
	case "version":
		fmt.Printf("purtypics v%s
", version)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s
", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Purtypics - Static Photo Gallery Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  purtypics <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  generate    Generate static gallery from photos")
	fmt.Println("  edit        Launch web-based metadata editor")
	fmt.Println("  version     Show version")
	fmt.Println("  help        Show this help")
	fmt.Println()
	fmt.Println("Run 'purtypics <command> -h' for command-specific options")
}

func runGenerate() {
	var (
		sourcePath   = flag.String("source", ".", "Source directory containing photo albums")
		outputPath   = flag.String("output", "", "Output directory (required)")
		baseURL      = flag.String("baseurl", "", "Base URL for the site")
		siteTitle    = flag.String("title", "Photo Gallery", "Site title")
		metadataPath = flag.String("metadata", "", "Path to metadata file (default: source/gallery.yaml)")
		verbose      = flag.Bool("verbose", false, "Verbose output")
	)

	// Output directory is required
	if *outputPath == "" {
		fmt.Println("Error: Output directory is required")
		fmt.Println("Please specify an output directory with -output")
		fmt.Println("
Example: purtypics generate -source /path/to/photos -output /path/to/website")
		os.Exit(1)
	}

	// Convert paths to absolute
	absSource, err := filepath.Abs(*sourcePath)
	if err != nil {
		log.Fatalf("Error resolving source path: %v", err)
	}

	absOutput, err := filepath.Abs(*outputPath)
	if err != nil {
		log.Fatalf("Error resolving output path: %v", err)
	}
	
	// Prevent output directory from being the same as source
	if absSource == absOutput {
		fmt.Println("Error: Output directory cannot be the same as source directory")
		fmt.Println("
Using the source directory as output will clutter your original photos")
		fmt.Println("with thumbnails and resized files. Please choose a different output directory.")
		fmt.Println("
Example: purtypics generate -source /path/to/photos -output /path/to/website")
		os.Exit(1)
	}

	generator := gallery.NewGenerator(absSource, absOutput, *siteTitle, *baseURL, *verbose)
	generator.MetadataPath = *metadataPath
	
	if err := generator.Generate(); err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	fmt.Println("✨ Gallery generated successfully!")
}

func runEdit() {
	var (
		sourcePath   = flag.String("source", ".", "Source directory containing photos")
		outputPath   = flag.String("output", "", "Output directory for generated gallery (required)")
		metadataPath = flag.String("metadata", "", "Path to metadata file (default: source/gallery.yaml)")
		port         = flag.Int("port", 8080, "Port to run the editor on")
		noBrowser    = flag.Bool("no-browser", false, "Don't open browser automatically")
	)
	
	flag.Parse()

	// Verify source directory exists
	if _, err := os.Stat(*sourcePath); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", *sourcePath)
	}

	// Create and start the server
	server := editor.NewServer(*sourcePath, *metadataPath, *port)
	
	// Open browser if not disabled
	if !*noBrowser {
		url := fmt.Sprintf("http://localhost:%d", *port)
		go func() {
			// Wait a moment for server to start
			fmt.Println("Opening browser...")
			openBrowser(url)
		}()
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// openBrowser opens the URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}