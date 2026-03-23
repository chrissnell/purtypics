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
  - s3: Deploy to Amazon S3 with incremental sync
  - cloudflare: Deploy to Cloudflare Pages with hash-based dedup

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
		if err := deployer.Deploy(); err != nil {
			return err
		}
		fmt.Printf("\nDeployment completed successfully!\n")
		return nil
		
	case "s3":
		if config.S3 == nil {
			return fmt.Errorf("S3 configuration not found in deploy.yaml")
		}
		s3Deployer := deploy.NewS3Deployer(config.S3, outputPath)
		fmt.Printf("Deploying gallery to %s...\n", s3Deployer.GetInfo())
		if deployDryRun {
			fmt.Println("(DRY RUN - testing connection)")
			if err := s3Deployer.TestConnection(); err != nil {
				return err
			}
			fmt.Println("Connection test successful!")
		} else {
			if err := s3Deployer.Deploy(); err != nil {
				return err
			}
			fmt.Println("\nDeployment completed successfully!")
		}
		return nil
		
	case "cloudflare":
		if config.Cloudflare == nil {
			return fmt.Errorf("cloudflare configuration not found in deploy.yaml")
		}
		deployer, err := deploy.NewCloudflareDeployer(config.Cloudflare, outputPath)
		if err != nil {
			return err
		}
		fmt.Printf("Deploying gallery to %s...\n", deployer.GetInfo())
		if deployDryRun {
			fmt.Println("(DRY RUN - testing connection)")
			if err := deployer.TestConnection(); err != nil {
				return err
			}
			fmt.Println("Connection test successful!")
		} else {
			if err := deployer.Deploy(); err != nil {
				return err
			}
			fmt.Println("\nDeployment completed successfully!")
		}
		return nil
		
	default:
		return fmt.Errorf("unknown deployment target: %s", deployTarget)
	}
}