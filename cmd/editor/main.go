package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/cjs/purtypics/pkg/editor"
)

func main() {
	var (
		sourcePath   = flag.String("source", ".", "Source directory containing photos")
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
	
	// Start the server in a goroutine to get the actual port
	serverStarted := make(chan bool)
	go func() {
		serverStarted <- true
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	
	// Wait for server to bind to port
	<-serverStarted
	
	// Open browser if not disabled with actual port
	if !*noBrowser {
		url := fmt.Sprintf("http://localhost:%d", server.Port)
		fmt.Println("Opening browser...")
		openBrowser(url)
	}
	
	// Block forever
	select {}
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