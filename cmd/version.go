package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			fmt.Println("purtypics version: development")
		} else {
			fmt.Printf("purtypics version: %s\n", version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}