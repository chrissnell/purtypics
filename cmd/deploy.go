package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/cjs/purtypics/pkg/deploy"
)

var (
	deployDryRun bool
	deployTarget string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the generated gallery to a remote host",
	Long: `Deploy the generated gallery to a remote host using various methods.
	
Supported deployment methods:
  - rsync: Deploy to a remote server via SSH
  - s3: Deploy to Amazon S3 (coming soon)
  - cloudflare: Deploy to Cloudflare Pages (coming soon)

Configuration is read from deploy.yaml in the gallery directory.`,
	RunE: runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
	
	deployCmd.Flags().BoolVar(&deployDryRun, "dry-run", false, "Perform a dry run without making changes")
	deployCmd.Flags().StringVar(&deployTarget, "target", "", "Deployment target (rsync, s3, cloudflare)")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	// Get gallery path
	galleryPath := "."
	if len(args) > 0 {
		galleryPath = args[0]
	}

	// Load deployment config
	config, err := deploy.LoadConfig(galleryPath)
	if err != nil {
		return fmt.Errorf("loading deploy config: %w", err)
	}

	// Get output path
	outputPath := filepath.Join(galleryPath, "output")
	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("output directory not found. Please run 'purtypics generate' first")
	}

	// Determine deployment target
	if deployTarget == "" {
		// Auto-detect from config
		if config.Rsync != nil {
			deployTarget = "rsync"
		} else if config.S3 != nil {
			deployTarget = "s3"
		} else if config.Cloudflare != nil {
			deployTarget = "cloudflare"
		} else {
			return fmt.Errorf("no deployment configuration found. Please configure deploy.yaml or use 'purtypics edit' to set up deployment")
		}
	}

	// Apply dry run flag
	if deployDryRun {
		switch deployTarget {
		case "rsync":
			if config.Rsync != nil {
				config.Rsync.DryRun = true
			}
		}
	}

	// Execute deployment
	switch deployTarget {
	case "rsync":
		if config.Rsync == nil {
			return fmt.Errorf("rsync configuration not found in deploy.yaml")
		}
		deployer := deploy.NewRsyncDeployer(config.Rsync, outputPath)
		fmt.Printf("Deploying gallery via %s...\n", deployer.GetInfo())
		if deployDryRun {
			fmt.Println("(DRY RUN - no changes will be made)")
		}
		return deployer.Deploy()
		
	case "s3":
		return fmt.Errorf("S3 deployment not yet implemented")
		
	case "cloudflare":
		return fmt.Errorf("Cloudflare deployment not yet implemented")
		
	default:
		return fmt.Errorf("unknown deployment target: %s", deployTarget)
	}
}