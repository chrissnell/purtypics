package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string

var rootCmd = &cobra.Command{
	Use:   "purtypics",
	Short: "A photo gallery generator and metadata editor",
	Long: `purtypics is a tool for creating static photo galleries and editing photo metadata.
	
It provides two main commands:
  - generate: Creates a static gallery from your photos
  - edit: Launches a web-based metadata editor`,
}

func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}