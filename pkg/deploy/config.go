package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the deployment configuration
type Config struct {
	Rsync      *RsyncConfig      `yaml:"rsync,omitempty" json:"rsync,omitempty"`
	S3         *S3Config         `yaml:"s3,omitempty" json:"s3,omitempty"`
	Cloudflare *CloudflareConfig `yaml:"cloudflare,omitempty" json:"cloudflare,omitempty"`
}

// RsyncConfig represents rsync deployment configuration
type RsyncConfig struct {
	Host   string `yaml:"host" json:"host"`                // user@hostname
	Path   string `yaml:"path" json:"path"`                // remote path
	Port   int    `yaml:"port,omitempty" json:"port"`      // SSH port (default: 22)
	DryRun bool   `yaml:"dry_run,omitempty" json:"dry_run"` // perform dry run
}

// S3Config represents S3 deployment configuration
type S3Config struct {
	Bucket         string `yaml:"bucket" json:"bucket"`
	Region         string `yaml:"region" json:"region"`
	CloudFrontID   string `yaml:"cloudfront_id,omitempty" json:"cloudfront_id,omitempty"`
	StorageClass   string `yaml:"storage_class,omitempty" json:"storage_class,omitempty"`
	CacheControl   string `yaml:"cache_control,omitempty" json:"cache_control,omitempty"`
	ACL            string `yaml:"acl,omitempty" json:"acl,omitempty"`
}

// CloudflareConfig represents Cloudflare Pages deployment configuration
type CloudflareConfig struct {
	Project   string `yaml:"project" json:"project"`
	AccountID string `yaml:"account_id" json:"account_id"`
	Branch    string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// LoadConfig loads deployment configuration from deploy.yaml
func LoadConfig(galleryPath string) (*Config, error) {
	configPath := filepath.Join(galleryPath, "deploy.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty config if file doesn't exist
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading deploy config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing deploy config: %w", err)
	}

	// Set defaults
	if config.Rsync != nil && config.Rsync.Port == 0 {
		config.Rsync.Port = 22
	}

	return &config, nil
}

// SaveConfig saves deployment configuration to deploy.yaml
func SaveConfig(galleryPath string, config *Config) error {
	configPath := filepath.Join(galleryPath, "deploy.yaml")
	
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling deploy config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing deploy config: %w", err)
	}

	return nil
}