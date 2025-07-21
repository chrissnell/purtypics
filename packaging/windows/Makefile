# Windows Packaging Makefile for Purtypics

# Variables
BINARY_NAME = purtypics
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR = ../../build
DIST_DIR = ../../dist
MAIN_PACKAGE = ../..

# Go build settings
GOOS = windows
GOARCH = amd64
GOFLAGS = -v
LDFLAGS = -s -w -X main.version=$(VERSION)

# Default target
.DEFAULT_GOAL := build

# Build Windows binary
.PHONY: build
build: ## Build Windows executable
	@echo "Building Purtypics $(VERSION) for Windows AMD64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Binary created at: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# Create MSI installer using WiX (cross-platform)
.PHONY: msi
msi: build ## Create MSI installer using WiX
	@echo "Creating MSI installer..."
	@command -v wixl >/dev/null 2>&1 || { echo "wixl not found. Install msitools: apt-get install msitools (Debian/Ubuntu) or brew install msitools (macOS)"; exit 1; }
	@mkdir -p $(DIST_DIR)
	wixl -v purtypics.wxs -o $(DIST_DIR)/purtypics-$(VERSION)-windows-amd64.msi
	@echo "MSI installer created at: $(DIST_DIR)/purtypics-$(VERSION)-windows-amd64.msi"

# Create self-extracting archive (cross-platform alternative)
.PHONY: sfx
sfx: build ## Create self-extracting archive
	@echo "Creating self-extracting archive..."
	@mkdir -p $(DIST_DIR)
	@mkdir -p temp_sfx
	@cp $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe temp_sfx/
	@cp ../../README.md temp_sfx/
	@cp ../../LICENSE temp_sfx/ 2>/dev/null || echo "No LICENSE file found"
	@cd temp_sfx && zip -r ../$(DIST_DIR)/purtypics-$(VERSION)-windows-amd64.zip *
	@rm -rf temp_sfx
	@echo "Self-extracting archive created at: $(DIST_DIR)/purtypics-$(VERSION)-windows-amd64.zip"

# Create NSIS installer (requires makensis)
.PHONY: nsis
nsis: build ## Create NSIS installer (requires makensis)
	@echo "Creating NSIS installer..."
	@command -v makensis >/dev/null 2>&1 || { echo "makensis not found. Install NSIS: apt-get install nsis (Linux) or brew install nsis (macOS)"; exit 1; }
	@mkdir -p $(DIST_DIR)
	makensis -DVERSION=$(VERSION) purtypics.nsi
	@echo "NSIS installer created at: $(DIST_DIR)/purtypics-setup-$(VERSION)-windows-amd64.exe"

# Create all package types
.PHONY: all
all: build sfx ## Build binary and create packages

# Clean build artifacts
.PHONY: clean
clean: ## Clean build artifacts
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@rm -rf $(DIST_DIR)/purtypics-*-windows-amd64.*
	@rm -rf temp_sfx

# Help
.PHONY: help
help: ## Show this help message
	@echo "Windows Packaging Makefile for Purtypics"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m%-15s\033[0m %s\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-13s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)